package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"

    "github.com/charlie-haley/home-infra/template/internal/manifest"
    ksv1 "github.com/fluxcd/kustomize-controller/api/v1"
    fluxmetav1 "github.com/fluxcd/pkg/apis/meta"
    sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
    "github.com/urfave/cli/v2"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/cli-runtime/pkg/printers"
    "sigs.k8s.io/yaml"
)

var frameworkResources []runtime.Object

var outputDir string
var gitSha string
var image string
var registry string
var user string
var pass string
var publish bool

func main() {
    var dir string
    app := &cli.App{
        Name: "mirdain",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:        "dir",
                Value:       "manifests",
                Usage:       "directory of yaml manifests",
                Destination: &dir,
            },
            &cli.StringFlag{
                Name:        "output-dir",
                Value:       "output",
                Usage:       "output directory",
                Destination: &outputDir,
            },
            &cli.StringFlag{
                Name:        "git-sha",
                Value:       "dev",
                Usage:       "git commit SHA",
                EnvVars:     []string{"GIT_SHA"},
                Destination: &gitSha,
            },
            &cli.StringFlag{
                Name:        "image",
                Value:       "dev",
                Usage:       "git commit SHA",
                EnvVars:     []string{"IMAGE"},
                Destination: &image,
            },
            &cli.StringFlag{
                Name:        "registry",
                Value:       "ghcr.io",
                Usage:       "docker registry",
                EnvVars:     []string{"REGISTRY"},
                Destination: &registry,
            },
            &cli.StringFlag{
                Name:        "user",
                Usage:       "docker registry user",
                EnvVars:     []string{"USER"},
                Destination: &user,
            },
            &cli.StringFlag{
                Name:        "pass",
                Usage:       "docker registry pass",
                EnvVars:     []string{"PASS"},
                Destination: &pass,
            },
            &cli.BoolFlag{
                Name:        "publish",
                Value:       false,
                Usage:       "whether to publish the image",
                EnvVars:     []string{"PUBLISH"},
                Destination: &publish,
            },
        },

        Commands: []*cli.Command{
            {
                Name:    "template",
                Aliases: []string{"t"},
                Usage:   "template kubernetes manifests based on directory.",
                Action: func(ctx *cli.Context) error {
                    return template(dir)
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func template(dir string) error {
    manifests, _ := os.ReadDir(dir)

    // Clean output dir
    err := os.RemoveAll(outputDir)
    if err != nil {
        return err
    }

    // Create output dir
    err = os.Mkdir(outputDir, 0755)
    if err != nil {
        return err
    }

    for _, directory := range manifests {

        // Ignore non-directory entries
        if !directory.IsDir() {
            continue
        }

        namespace := directory.Name()
        fmt.Printf("🟧 Processing namespace %s...\n", namespace)
        err := os.Mkdir(outputDir+"/"+namespace+"/", 0755)
        if err != nil {
            return err
        }

        // Add namespace to framework output
        addNamespaceResource(namespace)

        // Read contents of namespace folder
        namespacePath := filepath.Join(dir, namespace)
        apps, _ := os.ReadDir(namespacePath)

        for _, file := range apps {
            app := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

            fmt.Printf("  🟧 Processing app %s...\n", app)
            err := os.Mkdir(outputDir+"/"+namespace+"/"+app, 0755)
            if err != nil {
                return err
            }

            appPath := filepath.Join(namespacePath, file.Name())
            f, _ := os.ReadFile(appPath)

            var m manifest.Manifest
            yaml.Unmarshal(f, &m)

            appDir := outputDir + "/" + namespace + "/" + app

            err = m.Process(app, namespace, appDir)
            if err != nil {
                return err
            }

            addKustomizationResource(app, namespace, m.DependsOn)
        }
    }

    err = createFrameworkResources()
    if err != nil {
        return err
    }

    // Push files, this should move to code at some point to remove dependency on Flux CLI
    if publish {
        // Docker login, as with the above, this can be moved to code at some point to reduce manual system calls
        out, err := exec.Command("bash", "-c", "echo" + pass + " | docker login" + registry + "-u" + user + "--password-stdin").CombinedOutput()
        println(string(out))
        if err != nil {
            return err
        }

        timestamp := time.Now()

        tag := fmt.Sprintf("%v-%v", gitSha[0:7], fmt.Sprint(timestamp.Unix()))
        frameworkFile := outputDir + "/framework.yaml"
        out, err = exec.Command("flux", "push", "artifact", "oci://"+image+"/kustomizations:"+tag, "--source", "https://github.com/charlie-haley/home-infra", "--revision", gitSha, "--path", frameworkFile).CombinedOutput()
        println(string(out))
        if err != nil {
            return err
        }
        os.Remove(frameworkFile)

        out, err = exec.Command("flux", "push", "artifact", "oci://"+image+"/manifests:"+tag, "--source", "https://github.com/charlie-haley/home-infra", "--revision", gitSha, "--path", outputDir).CombinedOutput()
        println(string(out))
        return err
    }
    return nil
}

func createFrameworkResources() error {
    file, err := os.Create("output/framework.yaml")
    if err != nil {
        return err
    }

    defer file.Close()

    y := printers.YAMLPrinter{}
    for _, obj := range frameworkResources {
        y.PrintObj(obj, file)
    }
    return nil
}

func addNamespaceResource(namespace string) {
    annotations := map[string]string{"volsync.backube/privileged-movers": "true"}
    if namespace == "data" || namespace == "media" || namespace == "storage" {
        annotations["pod-security.kubernetes.io/warn"] = "privileged"
        annotations["pod-security.kubernetes.io/enforce"] = "privileged"
    }
    ns := &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name:        namespace,
            Annotations: annotations,
        },
    }
    ns.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"})
    frameworkResources = append(frameworkResources, ns)
}

func addKustomizationResource(app string, namespace string, dependsOn []fluxmetav1.NamespacedObjectReference) {
    // Create Kustomization
    ks := &ksv1.Kustomization{
        ObjectMeta: metav1.ObjectMeta{
            Name:      namespace + "-" + app,
            Namespace: "flux-system",
        },
        Spec: ksv1.KustomizationSpec{
            Interval: metav1.Duration{
                Duration: 15 * time.Minute,
            },
            Wait:            true,
            TargetNamespace: namespace,
            Prune:           true,
            SourceRef: ksv1.CrossNamespaceSourceReference{
                Kind: sourcev1.OCIRepositoryKind,
                Name: "manifests",
            },
            Path:      "./" + namespace + "/" + app,
            DependsOn: dependsOn,
            PostBuild: &ksv1.PostBuild{
                SubstituteFrom: []ksv1.SubstituteReference{
                    {
                        Kind: "Secret",
                        Name: "vilya-flux",
                    },
                    {
                        Kind: "ConfigMap",
                        Name: "vilya-flux",
                    },
                },
            },
            Decryption: &ksv1.Decryption{
                Provider: "sops",
                SecretRef: &fluxmetav1.LocalObjectReference{
                    Name: "sops-age",
                },
            },
        },
    }
    ks.SetGroupVersionKind(schema.GroupVersionKind{Group: "kustomize.toolkit.fluxcd.io/v1", Version: "v1", Kind: "Kustomization"})
    frameworkResources = append(frameworkResources, ks)
}

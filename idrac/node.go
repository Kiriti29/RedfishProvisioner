package main


import (
    "fmt"
    "os"
    "log"
    "net/http"
    "io/ioutil"
    "strings"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/kubernetes"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    corev1 "k8s.io/api/core/v1"
)


type error interface {
    String() string
}

type node struct {
    gorm.Model
    UUID    string   `json:"UUID"`
    Name       string `json:"Name"`
    DeployStatus string `json:"DeployStatus"`
    ImageURL string `json:ImageURL`
}

func getDB() (*gorm.DB) {
    db, err := gorm.Open("sqlite3", "/opt/db/nodes.db")
    if err != nil {
      panic("failed to connect database")
    }
    return db

}

func Create(w http.ResponseWriter, r *http.Request) {

    db := getDB()
    defer db.Close()

    var newNode node
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
    }

    json.Unmarshal(reqBody, &newNode)
    errors := db.Create(&newNode).Error
    if errors != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"msg": "%s"}`, errors)

    } else {
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(newNode)
    }
}


func getNode(w http.ResponseWriter, r *http.Request) {

    db := getDB()
    defer db.Close()

    nodeUUID := mux.Vars(r)["uuid"]
    var singleNode node
    errors := db.First(&singleNode, "UUID = ?", nodeUUID).Error
    if errors != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"msg": "%s"}`, errors)

    } else {
        json.NewEncoder(w).Encode(singleNode)
    }


}

func updateNode(w http.ResponseWriter, r *http.Request) {

    db := getDB()
    defer db.Close()

    nodeUUID := mux.Vars(r)["uuid"]
    var updatedNode node

    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
    }
    json.Unmarshal(reqBody, &updatedNode)

    var existingNode node
    errors := db.First(&existingNode, "UUID = ?", nodeUUID).Error
    if errors != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"msg": "%s"}`, errors)
    }

    existingNode.Name = updatedNode.Name
    existingNode.DeployStatus = updatedNode.DeployStatus
    existingNode.ImageURL = updatedNode.ImageURL
    errors = db.Save(&existingNode).Error
    if errors != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"msg": "%s"}`, errors)
    } else {
        fmt.Fprintf(w, `{"msg": "Node (id:%d) updated successfully"}`, existingNode.ID)
    }

}

func deleteNode(w http.ResponseWriter, r *http.Request) {

    db := getDB()
    defer db.Close()

    nodeUUID := mux.Vars(r)["id"]
    var existingNode node
    errors := db.First(&existingNode, "UUID = ?", nodeUUID).Error
    if errors != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"msg": "%s"}`, errors)
    }
    errors = db.Delete(&existingNode).Error
    if errors != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"msg": "%s"}`, errors)
    } else {
        fmt.Fprintf(w, `{"msg": "Node (id:%d) deleted successfully"}`, existingNode.ID)
    }
}

func GetClusterConfig(){
    // creates the in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
      panic(err.Error())
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
      panic(err.Error())
    }
    return clientset
}

func CreateSecret(node_id string, namespace string, data []byte)  {
    clientset := GetClusterConfig()
    secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "apps/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      node_id + "-kubeconfig",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"config": data,
		},
		Type: "Opaque",
	}


	_, err = clientset.CoreV1().Secrets(namespace).Create(&secret)
	if err != nil {
        panic(err.Error())
	}
}

func saveKubeConfig(w http.ResponseWriter, r *http.Request) {

    kubeconfig_root_path := "/opt/kubeconfigs/"
    url_parts := strings.Split(r.URL.Path, "/")
    nodeUUID := url_parts[len(url_parts)-1]
    fmt.Println(nodeUUID)
    body, err := ioutil.ReadAll(r.Body)
    fmt.Println(body)
    if err != nil {
     panic(err)
    }
    data := []byte(body)
    os.MkdirAll(kubeconfig_root_path + nodeUUID, os.ModePerm)
    err = ioutil.WriteFile(kubeconfig_root_path + nodeUUID + "/kube_config", data, 0644)
    if err != nil {
     panic(err)
    }
    CreateSecret(nodeUUID, "metalkube", data)
    fmt.Fprintf(w, `{"msg": "Success", "node": "%s" }`, nodeUUID)
}



func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome Home")
}


func main() {

    db := getDB()
    defer db.Close()
    listenPort := ":" + os.Getenv("PROVISIONER_PORT")

    db.AutoMigrate(&node{})


    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homePage)
    router.HandleFunc("/node", Create).Methods("POST")
    router.HandleFunc("/nodes/{uuid}", getNode).Methods("GET")
    router.HandleFunc("/nodes/{uuid}", updateNode).Methods("POST")
    router.HandleFunc("/nodes/{uuid}", deleteNode).Methods("DELETE")
    router.HandleFunc("/update_kube_config/{uuid}", saveKubeConfig).Methods("POST")

    log.Fatal(http.ListenAndServe(listenPort, router))


}

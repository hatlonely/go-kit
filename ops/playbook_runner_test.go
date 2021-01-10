package ops

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTopOrder(t *testing.T) {
	Convey("TestTopOrder", t, func() {
		m := map[string]string{
			"NAMESPACE":         "prod",
			"NAME":              "rpc-cicd",
			"IMAGE_REPOSITORY":  "${NAME}",
			"IMAGE_TAG":         "$(echo ${CLUSTER_NAME_1})",
			"CLUSTER_NAME_1":    "XX_${NAME}_${REGION_ID}_YY",
			"CLUSTER_NAME_2":    "${NAME}_${REGION_ID}",
			"CLUSTER_NAME_3":    "${CLUSTER_NAME_2}_${NAME}_${REGION_ID}",
			"REGION_ID":         "cn-shanghai",
			"REGISTRY_USERNAME": "{{.registry.username}}",
			"REGISTRY_PASSWORD": "{{.registry.password}}",
		}
		order, err := topOrderEnvMap(m)
		So(err, ShouldBeNil)
		So(len(order), ShouldEqual, 10)
		for _, key := range order {
			_, ok := m[key]
			So(ok, ShouldBeTrue)
		}
		fmt.Println(order)
	})
}

func TestParseEnvironment(t *testing.T) {
	Convey("TestParseEnvironment", t, func() {
		environment, err := ParseEnvironment(map[string]map[string]string{
			"default": {
				"NAMESPACE":        "prod",
				"NAME":             "rpc-cicd",
				"IMAGE_REPOSITORY": "${NAME}",
				"IMAGE_TAG":        "$(echo hello world)",
				"CLUSTER_NAME_1":   "XX_${NAME}_${REGION_ID}_YY",
				"CLUSTER_NAME_2":   "${NAME}_${REGION_ID}",
				"CLUSTER_NAME_3":   "${CLUSTER_NAME_2}_${NAME}_${REGION_ID}",
			},
			"prod": {
				"NAME":              "rpc-cicd-1",
				"REGION_ID":         "cn-shanghai",
				"REGISTRY_USERNAME": "{{.registry.username}}",
				"REGISTRY_PASSWORD": "{{.registry.password}}",
				"NAME1":             "${NAME//-/_}",
			},
		}, "tmp", "prod")

		So(err, ShouldBeNil)
		So(environment, ShouldResemble, []string{
			"CLUSTER_NAME_1=XX_rpc-cicd-1_cn-shanghai_YY",
			"CLUSTER_NAME_2=rpc-cicd-1_cn-shanghai",
			"CLUSTER_NAME_3=rpc-cicd-1_cn-shanghai_rpc-cicd-1_cn-shanghai",
			"DEP=tmp/dep",
			"ENV=prod",
			"IMAGE_REPOSITORY=rpc-cicd-1",
			"IMAGE_TAG=hello world",
			"NAME1=rpc_cicd_1",
			"NAME=rpc-cicd-1",
			"NAMESPACE=prod",
			"REGION_ID=cn-shanghai",
			"REGISTRY_PASSWORD={{.registry.password}}",
			"REGISTRY_USERNAME={{.registry.username}}",
			"TMP=tmp/env/prod",
		})
	})
}

func TestExecCommand(t *testing.T) {
	Convey("TestExecCommand", t, func() {
		status, stdout, stderr, err := ExecCommand("echo ${EXEC_COMMAND_KEY}=${EXEC_COMMAND_VALUE}", []string{
			"EXEC_COMMAND_KEY=key", "EXEC_COMMAND_VALUE=val",
		}, "")
		So(status, ShouldEqual, 0)
		So(stdout, ShouldEqual, "key=val\n")
		So(stderr, ShouldEqual, "")
		So(err, ShouldBeNil)
	})
}

func TestNewCICDRunner(t *testing.T) {
	Convey("TestNewCICDRunner", t, func() {
		os.MkdirAll("tmp", 0755)
		ioutil.WriteFile(`tmp/test.yaml`, []byte(`name: rpc-cicd
env:
  default:
    NAMESPACE: "prod"
    NAME: "rpc-cicd"
    ELASTICSEARCH_SERVER: "elasticsearch-master:9200"
    IMAGE_PULL_SECRET: "hatlonely-pull-secrets"
    IMAGE_REPOSITORY: "${NAME}"
    IMAGE_TAG: "$(cd .. && git describe --tags | awk '{print(substr($0,2,length($0)))}')"
    REPLICA_COUNT: 3
    INGRESS_HOST: "k8s.cicd.hatlonely.com"
    INGRESS_SECRET: "k8s-secret"
    K8S_CONTEXT: "homek8s"
  prod:
    REGISTRY_SERVER: "{{.registry.server}}"
    REGISTRY_USERNAME: "{{.registry.username}}"
    REGISTRY_PASSWORD: "{{.registry.password}}"
    REGISTRY_NAMESPACE: "{{.registry.namespace}}"
    MONGO_URI: "mongodb://{{.mongo.username}}:{{.mongo.password}}@{{.mongo.server}}"

task:
  image:
    args:
      key1:
        default: val1
    step:
      - make test
      - make image
  test:
    step:
      - make test
`), 0644)
		ioutil.WriteFile(`tmp/root.json`, []byte(`{
  "registry": {
	"server": "registry.cn-beijing.aliyuncs.com",
	"username": "hatlonely@foxmail.com",
	"password": "123456",
	"namespace": "hatlonely"
  },
  "mongo": {
	"server": "mongo-mongodb.prod.svc.cluster.local",
	"username": "root",
	"password": "RF7A7UVPET"
  },
}`), 0644)

		runner, err := NewPlaybookRunner(`tmp/test.yaml`, `tmp/root.json`)
		So(err, ShouldBeNil)
		_, err = runner.Environment("prod")
		So(err, ShouldBeNil)
		So(runner.environment["prod"], ShouldResemble, []string{
			`DEP=tmp/dep`,
			`ELASTICSEARCH_SERVER=elasticsearch-master:9200`,
			"ENV=prod",
			`IMAGE_PULL_SECRET=hatlonely-pull-secrets`,
			`IMAGE_REPOSITORY=rpc-cicd`,
			runner.environment["prod"][5],
			`INGRESS_HOST=k8s.cicd.hatlonely.com`,
			`INGRESS_SECRET=k8s-secret`,
			`K8S_CONTEXT=homek8s`,
			`MONGO_URI=mongodb://root:RF7A7UVPET@mongo-mongodb.prod.svc.cluster.local`,
			`NAME=rpc-cicd`,
			`NAMESPACE=prod`,
			`REGISTRY_NAMESPACE=hatlonely`,
			`REGISTRY_PASSWORD=123456`,
			`REGISTRY_SERVER=registry.cn-beijing.aliyuncs.com`,
			`REGISTRY_USERNAME=hatlonely@foxmail.com`,
			`REPLICA_COUNT=3`,
			"TMP=tmp/env/prod",
		})
		So(runner.playbook.Task, ShouldResemble, map[string]struct {
			Args map[string]struct {
				Validation string
				Type       string `dft:"string"`
				Default    string
			}
			Step []string
		}{
			"image": {
				Args: map[string]struct {
					Validation string
					Type       string `dft:"string"`
					Default    string
				}{
					"key1": {
						Type:    "string",
						Default: "val1",
					},
				},
				Step: []string{
					"make test",
					"make image",
				},
			},
			"test": {
				Step: []string{
					"make test",
				},
			},
		})

		os.RemoveAll("tmp")
	})
}

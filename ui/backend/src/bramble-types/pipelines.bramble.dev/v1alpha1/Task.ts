import { IObjectMeta } from "@kubernetes-models/apimachinery/apis/meta/v1/ObjectMeta";
import { addSchema } from "@kubernetes-models/apimachinery/_schemas/IoK8sApimachineryPkgApisMetaV1ObjectMeta";
import { Model, setSchema, ModelData, createTypeMetaGuard } from "@kubernetes-models/base";
import { register } from "@kubernetes-models/validate";

const schemaId = "pipelines.bramble.dev.v1alpha1.Task";
const schema = {
  "type": "object",
  "properties": {
    "apiVersion": {
      "type": "string",
      "enum": [
        "pipelines.bramble.dev/v1alpha1"
      ]
    },
    "kind": {
      "type": "string",
      "enum": [
        "Task"
      ]
    },
    "metadata": {
      "oneOf": [
        {
          "$ref": "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta#"
        },
        {
          "type": "null"
        }
      ]
    },
    "spec": {
      "properties": {
        "command": {
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "dependencies": {
          "items": {
            "type": "string"
          },
          "type": "array",
          "nullable": true
        },
        "image": {
          "type": "string"
        }
      },
      "required": [
        "command",
        "image"
      ],
      "type": "object",
      "nullable": true
    },
    "status": {
      "type": "object",
      "properties": {},
      "nullable": true
    }
  },
  "required": [
    "apiVersion",
    "kind"
  ]
};

/**
 * Task is the Schema for the tasks API
 */
export interface ITask {
  /**
   * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
   */
  "apiVersion": "pipelines.bramble.dev/v1alpha1";
  /**
   * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
   */
  "kind": "Task";
  "metadata"?: IObjectMeta;
  /**
   * TaskSpec defines the desired state of Task
   */
  "spec"?: {
    /**
     * Command executed by the container, can be used to determine the behaviour of a CLI app.
     */
    "command": Array<string>;
    "dependencies"?: Array<string>;
    /**
     * Docker image which will be used.
     */
    "image": string;
  };
  /**
   * TaskStatus defines the observed state of Task
   */
  "status"?: {
  };
}

/**
 * Task is the Schema for the tasks API
 */
export class Task extends Model<ITask> implements ITask {
  "apiVersion": ITask["apiVersion"];
  "kind": ITask["kind"];
  "metadata"?: ITask["metadata"];
  "spec"?: ITask["spec"];
  "status"?: ITask["status"];

static apiVersion: ITask["apiVersion"] = "pipelines.bramble.dev/v1alpha1";
static kind: ITask["kind"] = "Task";
static is = createTypeMetaGuard<ITask>(Task);

constructor(data?: ModelData<ITask>) {
  super({
    apiVersion: Task.apiVersion,
    kind: Task.kind,
    ...data
  } as ITask);
}
}


setSchema(Task, schemaId, () => {
  addSchema();
  register(schemaId, schema);
});

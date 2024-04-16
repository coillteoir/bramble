import { IObjectMeta } from "@kubernetes-models/apimachinery/apis/meta/v1/ObjectMeta";
import { addSchema } from "@kubernetes-models/apimachinery/_schemas/IoK8sApimachineryPkgApisMetaV1ObjectMeta";
import { Model, setSchema, ModelData, createTypeMetaGuard } from "@kubernetes-models/base";
import { register } from "@kubernetes-models/validate";

const schemaId = "pipelines.bramble.dev.v1alpha1.Execution";
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
        "Execution"
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
        "branch": {
          "type": "string"
        },
        "pipeline": {
          "type": "string"
        },
        "repo": {
          "type": "string"
        }
      },
      "required": [
        "branch",
        "pipeline",
        "repo"
      ],
      "type": "object",
      "nullable": true
    },
    "status": {
      "properties": {
        "completedTasks": {
          "items": {
            "type": "string"
          },
          "type": "array",
          "nullable": true
        },
        "executing": {
          "items": {
            "type": "string"
          },
          "type": "array",
          "nullable": true
        },
        "phase": {
          "type": "string",
          "nullable": true
        },
        "repoCloned": {
          "type": "boolean",
          "nullable": true
        },
        "volumeProvisioned": {
          "type": "boolean",
          "nullable": true
        }
      },
      "type": "object",
      "nullable": true
    }
  },
  "required": [
    "apiVersion",
    "kind"
  ]
};

/**
 * Execution is the Schema for the executions API.
 */
export interface IExecution {
  /**
   * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
   */
  "apiVersion": "pipelines.bramble.dev/v1alpha1";
  /**
   * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
   */
  "kind": "Execution";
  "metadata"?: IObjectMeta;
  /**
   * ExecutionSpec defines the desired state of Execution.
   */
  "spec"?: {
    /**
     * Git branch.
     */
    "branch": string;
    /**
     * Reference to the pipeline which will be executed.
     */
    "pipeline": string;
    /**
     * Git repo hosting the code to be tested against pipeline.
     */
    "repo": string;
  };
  /**
   * ExecutionStatus defines the observed state of Execution.
   */
  "status"?: {
    /**
     * Tasks which have already succeeded.
     */
    "completedTasks"?: Array<string>;
    /**
     * Describes which tasks are currently running.
     */
    "executing"?: Array<string>;
    /**
     * Describes the state of the execution
     */
    "phase"?: string;
    /**
     * Tells the controller if the repo is cloned.
     */
    "repoCloned"?: boolean;
    /**
     * Shows if the PV for this execution has been provisioned.
     */
    "volumeProvisioned"?: boolean;
  };
}

/**
 * Execution is the Schema for the executions API.
 */
export class Execution extends Model<IExecution> implements IExecution {
  "apiVersion": IExecution["apiVersion"];
  "kind": IExecution["kind"];
  "metadata"?: IExecution["metadata"];
  "spec"?: IExecution["spec"];
  "status"?: IExecution["status"];

static apiVersion: IExecution["apiVersion"] = "pipelines.bramble.dev/v1alpha1";
static kind: IExecution["kind"] = "Execution";
static is = createTypeMetaGuard<IExecution>(Execution);

constructor(data?: ModelData<IExecution>) {
  super({
    apiVersion: Execution.apiVersion,
    kind: Execution.kind,
    ...data
  } as IExecution);
}
}


setSchema(Execution, schemaId, () => {
  addSchema();
  register(schemaId, schema);
});

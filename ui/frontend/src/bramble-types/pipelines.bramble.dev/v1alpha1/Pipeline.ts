import { IObjectMeta } from "@kubernetes-models/apimachinery/apis/meta/v1/ObjectMeta";
import { addSchema } from "@kubernetes-models/apimachinery/_schemas/IoK8sApimachineryPkgApisMetaV1ObjectMeta";
import { Model, setSchema, ModelData, createTypeMetaGuard } from "@kubernetes-models/base";
import { register } from "@kubernetes-models/validate";

const schemaId = "pipelines.bramble.dev.v1alpha1.Pipeline";
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
        "Pipeline"
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
        "taskRefs": {
          "items": {
            "properties": {
              "dependencies": {
                "items": {
                  "type": "string"
                },
                "type": "array",
                "nullable": true
              },
              "name": {
                "type": "string"
              }
            },
            "required": [
              "name"
            ],
            "type": "object"
          },
          "type": "array",
          "nullable": true
        },
        "tasks": {
          "items": {
            "properties": {
              "name": {
                "type": "string"
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
                  },
                  "workdir": {
                    "type": "string",
                    "nullable": true

                  }
                },
                "required": [
                  "command",
                  "image"
                ],
                "type": "object"
              }
            },
            "required": [
              "name",
              "spec"
            ],
            "type": "object"
          },
          "type": "array",
          "nullable": true
        }
      },
      "type": "object",
      "nullable": true
    },
    "status": {
      "properties": {
        "taskscreated": {
          "type": "boolean",
          "nullable": true
        },
        "validdeps": {
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
 * Pipeline is the Schema for the pipelines API
 */
export interface IPipeline {
  /**
   * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
   */
  "apiVersion": "pipelines.bramble.dev/v1alpha1";
  /**
   * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
   */
  "kind": "Pipeline";
  "metadata"?: IObjectMeta;
  /**
   * PipelineSpec defines the desired state of Pipeline
   */
  "spec"?: {
    /**
     * Allows developers to use pre applied tasks in the same namespace
     */
    "taskRefs"?: Array<{
      "dependencies"?: Array<string>;
      "name": string;
    }>;
    /**
     * Allows developers to create a list of tasks
     */
    "tasks"?: Array<{
      /**
       * Name of task to be ran
       */
      "name": string;
      /**
       * Spec of given task
       */
      "spec": {
        /**
         * Command executed by the container, can be used to determine the behaviour of a CLI app.
         */
        "command": Array<string>;
        /**
         * Tasks which need to run before current task
         */
        "dependencies"?: Array<string>;
        /**
         * Docker image which will be used.
         */
        "image": string;
        /**
         * Directory in which task container executes
         */
        "workdir"?: string;
      };
    }>;
  };
  /**
   * PipelineStatus defines the observed state of Pipeline
   */
  "status"?: {
    "taskscreated"?: boolean;
    /**
     * Shows that the pipeline has a valid dependency tree.
     */
    "validdeps"?: boolean;
  };
}

/**
 * Pipeline is the Schema for the pipelines API
 */
export class Pipeline extends Model<IPipeline> implements IPipeline {
  "apiVersion": IPipeline["apiVersion"];
  "kind": IPipeline["kind"];
  "metadata"?: IPipeline["metadata"];
  "spec"?: IPipeline["spec"];
  "status"?: IPipeline["status"];

static apiVersion: IPipeline["apiVersion"] = "pipelines.bramble.dev/v1alpha1";
static kind: IPipeline["kind"] = "Pipeline";
static is = createTypeMetaGuard<IPipeline>(Pipeline);

constructor(data?: ModelData<IPipeline>) {
  super({
    apiVersion: Pipeline.apiVersion,
    kind: Pipeline.kind,
    ...data
  } as IPipeline);
}
}
}



setSchema(Pipeline, schemaId, () => {
  addSchema();
  register(schemaId, schema);
});

import { Pipeline, PLtask } from "./bramble_types";

const k8s = require("@kubernetes/client-node");

const kc = new k8s.KubeConfig();

if (process.env.IN_CLUSTER === "1") {
  kc.loadFromCluster();
  console.log("Connecting to cluster from within");
} else {
  kc.loadFromDefault();
  console.log("Connecting to cluster via kubeconfig");
}

const k8sCRDApi = kc.makeApiClient(k8s.CustomObjectsApi);
const k8sCoreApi = kc.makeApiClient(k8s.CoreV1Api);

export const getPods = async (namespace: string) => {
  try {
    const response = await k8sCoreApi.listNamespacedPod(namespace);
    return response?.body?.items;
  } catch (err) {
    console.error(err);
  }
};

export const getPipelines = async (ns: string): Promise<Pipeline[] | Error> => {
  try {
    const response = await k8sCRDApi.listNamespacedCustomObject(
      "pipelines.bramble.dev",
      "v1alpha1",
      ns,
      "pipelines",
    );
    const pipelines: Pipeline[] = response?.body?.items.map((pipeline: any) => {
      return new Pipeline(
        { name: pipeline.metadata.name, namespace: ns },
        {
          tasks: pipeline.spec.tasks.map((task: PLtask) => {
            return new PLtask(task.name, task.spec);
          }),
        },
      );
    });
    return pipelines;
  } catch (err: any) {
    const ret = new Error(err.message);
    throw ret;
    return ret;
  }
};

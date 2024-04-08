import { pipelinesBrambleDev } from "./bramble-types";

const k8s = require("@kubernetes/client-node");

const kc = new k8s.KubeConfig();

if (process.env.KUBERNETES_SERVICE_HOST !== "") {
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

export const getPipelines = async (
  namespace: string,
): Promise<pipelinesBrambleDev.v1alpha1.Pipeline[] | Error> => {
  try {
    const response: any = await k8sCRDApi.listNamespacedCustomObject(
      "pipelines.bramble.dev",
      "v1alpha1",
      namespace,
      "pipelines",
    );
    const pipelines: pipelinesBrambleDev.v1alpha1.Pipeline[] =
      response?.body?.items.map(
        (pipeline: any) => new pipelinesBrambleDev.v1alpha1.Pipeline(pipeline),
      );
    return pipelines;
  } catch (err: any) {
    const ret = new Error(err.message);
    throw ret;
  }
};

export const getExecutions = async (
  namespace: string,
): Promise<pipelinesBrambleDev.v1alpha1.Execution[] | Error> => {
  try {
    const response: any = await k8sCRDApi.listNamespacedCustomObject(
      "pipelines.bramble.dev",
      "v1alpha1",
      namespace,
      "executions",
    );
    const executions: pipelinesBrambleDev.v1alpha1.Execution[] =
      response?.body?.items.map(
        (execution: any) =>
          new pipelinesBrambleDev.v1alpha1.Execution(execution),
      );
    return executions;
  } catch (err: any) {
    const ret = new Error(err.message);
    throw ret;
  }
};

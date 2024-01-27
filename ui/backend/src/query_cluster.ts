import { pipelinesBrambleDev } from "./bramble-types";

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

export const getPipelines = async (
  namespace: string,
): Promise<pipelinesBrambleDev.v1alpha1.Pipeline[] | Error> => {
  try {
    const response = await k8sCRDApi.listNamespacedCustomObject(
      "pipelines.bramble.dev",
      "v1alpha1",
      namespace,
      "pipelines",
    );
    console.log(response)
    const pls: pipelinesBrambleDev.v1alpha1.Pipeline[] =
      response?.body?.items.map((pl: pipelinesBrambleDev.v1alpha1.Pipeline) => {
        return new pipelinesBrambleDev.v1alpha1.Pipeline(pl);
      });
    return pls;
  } catch (err: any) {
    const ret = new Error(err.message);
    throw ret;
    return ret;
  }
};

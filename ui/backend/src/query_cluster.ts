import { Pod } from "kubernetes-types/core/v1";
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

export const getPL = async (ns: string): Promise<Pipeline[] | Error> => {
  try {
    const response = await k8sCRDApi.listNamespacedCustomObject(
      "pipelines.bramble.dev",
      "v1alpha1",
      ns,
      "pipelines",
    );
    const pls: Pipeline[] = response?.body?.items.map((pl: any) => {
      return new Pipeline(
        { name: pl.metadata.name, namespace: ns },
        {
          tasks: pl.spec.tasks.map((task: any) => {
            return new PLtask(task.name, task.spec, task?.dependencies);
          }),
        },
      );
    });
    return pls;
  } catch (err: any) {
    const ret = new Error(err.message);
    throw ret;
    return ret;
  }
};

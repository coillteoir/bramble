import { Pod } from "kubernetes-types/core/v1";
import * as bramble_types from "./bramble_types";

const k8s = require("@kubernetes/client-node");

const kc = new k8s.KubeConfig();
kc.loadFromDefault();

const k8sApi = kc.makeApiClient(k8s.CoreV1Api);
const k8sCRDApi = kc.makeApiClient(k8s.CustomObjectsApi);

export const getPo = async (ns: string) => {
  const response = await k8sApi.listNamespacedPod(ns);
  const podNames = response.body.items.map((pod: Pod) => ({
    name: pod?.metadata?.name,
    image: pod?.spec?.containers[0].image,
  }));
  return podNames;
};

export const getPL = async (ns: string): Promise<bramble_types.Pipeline[]> => {
  const response = await k8sCRDApi.listNamespacedCustomObject(
    "pipelines.bramble.dev",
    "v1alpha1",
    ns,
    "pipelines",
  );
  const pls: bramble_types.Pipeline[] = response?.body?.items.map((pl: any) => {
    console.log(pl)
    new bramble_types.Pipeline(
      { name: pl.metadata.name, namespace: ns },
      {
        tasks: pl.spec.tasks.map((task: any) => {
          new bramble_types.PLtask(task.name, task.spec, task?.dependencies);
        }),
      },
    );
  });
  return pls;
};

import { pipelinesBrambleDev } from "./bramble-types";
import { getLogger } from "log4js";
const k8s = require("@kubernetes/client-node");

const logger = getLogger();
logger.level = "info";

const kc = new k8s.KubeConfig();

if (process.env.KUBERNETES_SERVICE_HOST) {
  logger.info("Connecting to cluster from within");
  try {
    kc.loadFromCluster();
  } catch (err: any) {
    logger.error(err);
  }
  logger.info("CONNECTED!");
} else {
  logger.info("Connecting to cluster via kubeconfig");
  try {
    kc.loadFromDefault();
  } catch (err: any) {
    logger.error(err);
  }
  logger.info("CONNECTED!");
}

const k8sCRDApi = kc.makeApiClient(k8s.CustomObjectsApi);
const k8sBatchApi = kc.makeApiClient(k8s.BatchV1Api);

export const getJobs = async (namespace: string) => {
  try {
    const response = await k8sBatchApi.listNamespacedJob(namespace);
    return response?.body?.items;
  } catch (err: any) {
    const ret = new Error(err.message);
    logger.error(ret);
    throw ret;
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
    logger.error(ret);
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
    logger.error(ret);
    throw ret;
  }
};

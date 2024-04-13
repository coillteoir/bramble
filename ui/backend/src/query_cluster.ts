import { getLogger } from "log4js";

import { Job } from "kubernetes-models/batch/v1";
import { Pod } from "kubernetes-models/v1";
import { pipelinesBrambleDev } from "./bramble-types";
import {getLogger} from "log4js"
const k8s = require("@kubernetes/client-node");

import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;

const logger = getLogger();
logger.level = "info";

const kc = new k8s.KubeConfig();

if (process.env.KUBERNETES_SERVICE_HOST) {
    logger.info("Connecting to cluster from within");
    try {
        kc.loadFromCluster();
    } catch (err: any) {
        logger.error(err);
        process.exit(1);
    }
} else {
    logger.info("Connecting to cluster via kubeconfig");
    try {
        kc.loadFromDefault();
    } catch (err: any) {
        logger.error(err);
        process.exit(1);
    }

}

logger.info("CONNECTED!");

const k8sCoreApi = kc.makeApiClient(k8s.CoreV1Api);
const k8sBatchApi = kc.makeApiClient(k8s.BatchV1Api);

const k8sCRDApi = kc.makeApiClient(k8s.CustomObjectsApi);

export const getNamespacedJobLogs = async (
    namespace: string,
    jobName: string
) => {
    logger.info(`querying jobs in ${namespace} for logs`);
    try {
        const jobres = await getJobs(namespace);
        logger.info(`querying pods in ${namespace} for logs`);
        try {

            const logJob: Job | undefined = jobres.body?.items?.find(
                (job: Job) => job.metadata?.name === jobName
            );
            if (!logJob) {
                return;
            }
            logger.info(`job ${jobName} found`);

            const podres = await k8sCoreApi.listNamespacedPod(namespace);
            logger.info(`querying for pods in job ${logJob.metadata?.name}`);
            const logPod: Pod | undefined = podres.body?.items?.find(
                (pod: Pod) =>
                    pod.metadata?.labels?.["job-name"] === logJob.metadata?.name
            );

            if (!logPod || !logPod.spec || !logPod.metadata) {
                return;
            }
            logger.info(
                `retrieving logs from container ${logPod.spec.containers[0]} in pod ${logPod.metadata.name}`
            );

            try {
                const logs = await k8sCoreApi.readNamespacedPodLog(
                    logPod.metadata?.name,
                    namespace
                );
                return logs;
            } catch (err: any) {
                logger.error(err);
                throw err;
            }
        } catch (err: any) {
            throw err;
        }
    } catch (err: any) {
        throw err;
    }
};

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
    namespace: string
): Promise<Pipeline[] | Error> => {
    try {
        const response: any = await k8sCRDApi.listNamespacedCustomObject(
            "pipelines.bramble.dev",
            "v1alpha1",
            namespace,
            "pipelines"
        );
        const pipelines: Pipeline[] = response?.body?.items.map(
            (pipeline: any) =>
                new pipelinesBrambleDev.v1alpha1.Pipeline(pipeline)
        );
        return pipelines;
    } catch (err: any) {
        const ret = new Error(err.message);
        logger.error(ret);
        throw ret;
    }
};

export const getExecutions = async (
    namespace: string
): Promise<Execution[] | Error> => {
    try {
        const response: any = await k8sCRDApi.listNamespacedCustomObject(
            "pipelines.bramble.dev",
            "v1alpha1",
            namespace,
            "executions"
        );
        const executions: Execution[] = response?.body?.items.map(
            (execution: any) =>
                new pipelinesBrambleDev.v1alpha1.Execution(execution)
        );
        return executions;
    } catch (err: any) {
        const ret = new Error(err.message);
        logger.error(ret);
        throw ret;
    }

};

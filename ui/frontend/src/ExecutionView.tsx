import React, { useState } from "react";

import ReactFlow, { Node, Edge } from "reactflow";
import { Job, IJobStatus } from "kubernetes-models/batch/v1";
import { pipelinesBrambleDev } from "./bramble-types";

import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;

import { getLayoutedElements } from "./Layout.tsx";
import { generateNodes, generateEdges } from "./ExecutionGraph.tsx";

const TaskView = (props: {
    task:
        | {
              name: string;
              spec: {
                  image: string;
                  command: string[];
                  dependencies?: string[] | undefined;
                  workdir?: string | undefined;
              };
          }
        | undefined;
    job: Job | undefined;
}) => {
    if (!props.task) return;
    const taskInfo = props.task;
    const jobStatus = props.job?.status;
    const getJobSeconds = (status: IJobStatus) => {
        if (!status.startTime) {
            return 0;
        }
        return Math.floor(
            status.completionTime
                ? Date.parse(status.completionTime) -
                      Date.parse(status.startTime)
                : Date.now() - Date.parse(status.startTime)
        );
    };
    return (
        <>
            <h1 className="text-2xl font-bold">Task: {taskInfo.name}</h1>
            <h2 className="text-xl font-bold">Image: {taskInfo.spec.image}</h2>
            <h2 className="text-lg font-bold">
                Command: <code>{taskInfo.spec.command.join("\n")}</code>
            </h2>
            {taskInfo.spec.dependencies && (
                <>
                    <h2 className="font-bold">Dependencies</h2>
                    <ul>
                        {taskInfo.spec.dependencies.map((dep, i) => (
                            <li key={i}>{dep}</li>
                        ))}
                    </ul>
                </>
            )}
            {taskInfo.spec.workdir && (
                <h2 className="font-bold">WorkDir: {taskInfo.spec.workdir}</h2>
            )}
            {jobStatus && (
                <p>
                    Duration:
                    {new Date(getJobSeconds(jobStatus))
                        .toISOString()
                        .slice(11, 19)}
                </p>
            )}
        </>
    );
};

const ExecutionView = (props: {
    pipeline: Pipeline;
    execution: Execution | undefined;
}): React.ReactNode => {
    const [jobs, setJobs] = useState<Array<Job>>(new Array<Job>());

    const [focusedTask, setFocusedTask] = useState<string>("manifests");

    const fetchNamespacedJobs = async (namespace: string | undefined) => {
        if (namespace === undefined) {
            console.log("undefined namespace");
            return;
        }

        console.log(`Fetching jobs in ${namespace} ${Date()}`);
        try {
            await fetch(`http://localhost:5555/resources/jobs/${namespace}`)
                .then((response) => response.json())
                .then((jsonData) => {
                    setJobs(jsonData);
                });
        } catch (error) {
            console.error(error);
        }
    };

    const layouted = getLayoutedElements(
        props.pipeline.spec?.tasks
            ? generateNodes(
                  props.pipeline.spec?.tasks,
                  jobs,
                  props.execution,
                  setFocusedTask
              )
            : ([] as Node[]),
        props.pipeline.spec?.tasks
            ? generateEdges(props.pipeline.spec?.tasks)
            : ([] as Edge[])
    );

    const nodes = [...layouted.nodes];
    const edges = [...layouted.edges];

    React.useEffect(() => {
        const interval = setInterval(
            () => fetchNamespacedJobs(props.execution?.metadata?.namespace),
            1000
        );
        return () => {
            clearInterval(interval);
        };
    });

    return (
        <div className="flex h-full flex-col">
            <h1 className="w-full bg-green-800 text-3xl font-bold">
                Pipeline: {props.pipeline.metadata?.name}
            </h1>
            {props.execution && (
                <h2 className="w-full bg-green-900 text-2xl font-bold">
                    Execution: {props.execution.metadata?.name}
                </h2>
            )}
            <div
                className="flex"
                style={{
                    width: "100%",
                    height: "100%",
                }}
            >
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    fitView
                    defaultEdgeOptions={{
                        animated: true,
                        style: {
                            animationDirection: "reverse",
                            stroke: "green",
                        },
                    }}
                />
                <div className="inline-block h-full w-2/5">
                    <TaskView
                        task={props.pipeline.spec?.tasks?.find(
                            (task) => task.name === focusedTask
                        )}
                        job={
                            props.execution &&
                            jobs?.find(
                                (job: Job) =>
                                    props.execution?.metadata?.name ===
                                        job.metadata?.labels?.[
                                            "bramble-execution"
                                        ] &&
                                    job.metadata?.labels?.["bramble-task"] ===
                                        focusedTask
                            )
                        }
                    />
                </div>
            </div>
        </div>
    );
};

export { ExecutionView };

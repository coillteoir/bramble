import { Node, Edge } from "reactflow";
import { Job } from "kubernetes-models/batch/v1";
import { pipelinesBrambleDev } from "./bramble-types";

import Execution = pipelinesBrambleDev.v1alpha1.Execution;

export enum ExecutionPhase {
    Success,
    Running,
    Pending,
    Failure,
}

export const PipelineStatusIcon = (props: { phase: ExecutionPhase }) => {
    switch (props.phase) {
        case ExecutionPhase.Success:
            return (
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="size-6 shrink-0 stroke-current"
                    fill="none"
                    viewBox="0 0 24 24"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
            );
        case ExecutionPhase.Running:
            return (
                <span className="loading loading-spinner loading-sm text-blue-800" />
            );
        case ExecutionPhase.Failure:
            return (
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="size-6 shrink-0 stroke-current"
                    fill="none"
                    viewBox="0 0 24 24"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
            );
        case ExecutionPhase.Pending:
            return <span className="loading loading-dots loading-sm"></span>;
    }
};

const jobStatusIcon = (job: Job | undefined) => {
    if (job === undefined) {
        return <PipelineStatusIcon phase={ExecutionPhase.Pending} />;
    }
    if (job.status?.succeeded) {
        return <PipelineStatusIcon phase={ExecutionPhase.Success} />;
    }
    if (job.status?.failed) {
        return <PipelineStatusIcon phase={ExecutionPhase.Failure} />;
    }
    if (job.status?.active) {
        return <PipelineStatusIcon phase={ExecutionPhase.Running} />;
    }
};

export const generateNodes = (
    tasks: {
        name: string;
        spec: {
            image: string;
            command: string[];
            dependencies?: string[] | undefined;
            workDir?: string | undefined;
        };
    }[],
    jobs: Job[] | undefined,
    execution: Execution | undefined,
    taskSetter: any
): Node[] =>
    tasks.map(
        (task: {
            name: string;
            spec: {
                image: string;
                command: string[];
                dependencies?: string[] | undefined;
            };
        }): Node => {
            // get pod of current task
            const taskJob: Job | undefined =
                execution &&
                jobs?.find(
                    (job: Job) =>
                        execution.metadata?.name ===
                            job.metadata?.labels?.["bramble-execution"] &&
                        job.metadata?.labels?.["bramble-task"] === task.name
                );
            const spinner = execution && jobStatusIcon(taskJob);

            return {
                id: task.name,
                width: 120,
                height: 50,
                position: { x: 0, y: 0 },
                data: {
                    label: (
                        <div
                            className="group"
                            onClick={() => taskSetter(task.name)}
                        >
                            <p className="">{task.name}</p>
                            {spinner}
                        </div>
                    ),
                },
            };
        }
    );

export const generateEdges = (
    tasks: {
        name: string;
        spec: {
            image: string;
            command: string[];
            dependencies?: string[] | undefined;
        };
    }[]
): Edge[] =>
    tasks
        .map(
            (task: {
                name: string;
                spec: {
                    image: string;
                    command: string[];
                    dependencies?: string[] | undefined;
                };
            }): Edge[] =>
                task.spec.dependencies
                    ? task.spec.dependencies.map(
                          (dep: string, i: number): Edge => ({
                              id: task.name + i.toString(),
                              target: dep,
                              source: task.name,
                          })
                      )
                    : ([] as Edge[])
        )
        .flat();

import { Node, Edge } from "reactflow";
import { Job } from "kubernetes-models/batch/v1";
import { pipelinesBrambleDev } from "./bramble-types";

const jobStatusIcon = (job: Job | undefined) => {
    if (job === undefined) {
        return <span className="loading loading-dots loading-sm"></span>;
    }
    if (job.status?.succeeded !== 0) {
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
    }
    if (job.status?.active !== 0) {
        return (
            <span className="loading loading-spinner loading-sm text-blue-800" />
        );
    }
    if (job.status.failed !== 0) {
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
    }
};

export const generateNodes = (
    tasks: {
        name: string;
        spec: {
            image: string;
            command: string[];
            dependencies?: string[] | undefined;
        };
    }[],
    jobs: Job[] | undefined,
    execution: pipelinesBrambleDev.v1alpha1.Execution | undefined
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
                        <div className="group">
                            <p className="">{task.name}</p>
                            <p className="hidden group-hover:block">
                                Image: {task.spec.image}
                            </p>
                            <p className="hidden group-hover:block">
                                Command: {task.spec.command}
                            </p>
                            {task.spec.dependencies && (
                                <p className="hidden group-hover:block">
                                    Dependencies: {task.spec.dependencies}
                                </p>
                            )}
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

import { Node, Edge } from "reactflow";
import { Pod } from "kubernetes-models/v1";
import { pipelinesBrambleDev } from "./bramble-types";

export const generateNodes = (
    tasks: {
        name: string;
        spec: {
            image: string;
            command: string[];
            dependencies?: string[] | undefined;
        };
    }[],
    pods: Pod[],
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
            const taskPod: Pod | undefined =
                execution &&
                pods.find(
                    (pod: Pod) =>
                        execution.metadata?.name ===
                            pod.metadata?.labels?.["bramble-execution"] &&
                        pod.metadata?.labels?.["bramble-task"] === task.name
                );
            const colour =
                execution &&
                (() => {
                    if (taskPod === undefined) {
                        return "orange";
                    }
                    switch (taskPod.status?.phase) {
                        case "Succeeded":
                            return "green";
                        case "Running":
                            return "blue";
                        case "Failed":
                            return "red";
                    }
                    return "orange";
                })();

            return {
                id: task.name,
                width: 120,
                height: 50,
                position: { x: 0, y: 0 },
                data: {
                    label: (
                        <div>
                            <p>{task.name}</p>
                            {colour && (
                                <svg
                                    viewBox="0 0 2 2"
                                    xmlns="http://www.w3.org/2000/svg"
                                    style={{
                                        width: "20%",
                                        height: "20%",
                                    }}
                                >
                                    <circle cx="1" cy="1" r="1" fill={colour} />
                                </svg>
                            )}
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

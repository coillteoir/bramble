import { Node, Edge } from "reactflow";
import { Pod } from "kubernetes-models/v1";
import { pipelinesBrambleDev } from "./bramble-types";

// https://www.benmvp.com/blog/how-to-create-circle-svg-gradient-loading-spinner/
//const SVGPhase = () => (
//    <svg
//        width="20"
//        height="20"
//        viewBox="0 0 200 200"
//        color="#3f51b5"
//        fill="none"
//        xmlns="http://www.w3.org/2000/svg"
//        className="inline-block"
//    >
//        <defs>
//            <linearGradient id="spinner-secondHalf">
//                <stop offset="0%" stopOpacity="0" stopColor="currentColor" />
//                <stop
//                    offset="100%"
//                    stopOpacity="0.5"
//                    stopColor="currentColor"
//                />
//            </linearGradient>
//            <linearGradient id="spinner-firstHalf">
//                <stop offset="0%" stopOpacity="1" stopColor="currentColor" />
//                <stop
//                    offset="100%"
//                    stopOpacity="0.5"
//                    stopColor="currentColor"
//                />
//            </linearGradient>
//        </defs>
//
//        <g strokeWidth="8">
//            <path
//                stroke="url(#spinner-secondHalf)"
//                d="M 4 100 A 96 96 0 0 1 196 100"
//            />
//            <path
//                stroke="url(#spinner-firstHalf)"
//                d="M 196 100 A 96 96 0 0 1 4 100"
//            />
//
//            <path
//                stroke="currentColor"
//                strokeLinecap="round"
//                d="M 4 100 A 96 96 0 0 1 4 98"
//            />
//        </g>
//
//        <animateTransform
//            from="0 0 0"
//            to="360 0 0"
//            attributeName="transform"
//            type="rotate"
//            repeatCount="indefinite"
//            dur="1300ms"
//        />
//    </svg>
//);

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
                        return (
                            <span className="loading loading-dots loading-sm"></span>
                        );
                    }
                    switch (taskPod.status?.phase) {
                        case "Succeeded":
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
                        case "Running":
                            return (
                                <span className="loading loading-spinner loading-sm text-blue-800" />
                            );
                        case "Failed":
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
                    return (
                        <span className="loading loading-dots loading-sm"></span>
                    );
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
                            {colour}
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

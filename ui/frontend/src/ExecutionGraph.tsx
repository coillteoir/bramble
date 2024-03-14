import { Node, Edge } from "reactflow";
import { Pod } from "kubernetes-models/v1";
import Dagre from "@dagrejs/dagre";
import {pipelinesBrambleDev} from "./bramble-types"

// https://codesandbox.io/p/sandbox/romantic-bas-z2v5wm?file=%2FApp.js%3A63%2C51&utm_medium=sandpack
export const getLayoutedElements = (nodes: Node[], edges: Edge[]) => {
    const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));
    g.setGraph({ rankdir: "TB" });

    edges.forEach((edge: Edge) => g.setEdge(edge.source, edge.target));

    // any type used because the typing of setNode seems to be incorrect.
    // It behaves as expected when entire node is passed in
    nodes.forEach((node: Node) => g.setNode(node.id, node as any));

    Dagre.layout(g);

    return {
        nodes: nodes.map((node: Node) => {
            const { x, y } = g.node(node.id);
            return { ...node, position: { x, y } };
        }),
        edges,
    };
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
    pods: Pod[],
    execution: pipelinesBrambleDev.v1alpha1.Execution
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
            const taskPod:Pod = pods.filter(
                (pod: Pod) =>
                    execution.metadata?.name ==
                        pod.metadata?.labels?.["bramble-execution"] &&
                    pod.metadata?.labels?.["bramble-task"] == task.name
            )[0];
            const colour = (() => {
                if(taskPod === undefined) {
                    return "orange"
                }
                switch (taskPod.status?.phase) {
                    case "Succeeded":
                        return "green";
                        break;
                    case "Running":
                        return "blue";
                        break;
                    case "Failed":
                        return "red";
                        break;
                }
                return "orange"
            })()

            return {
                id: task.name,
                width: 120,
                height: 50,
                position: { x: 0, y: 0 },
                data: {
                    label: (
                        <div>
                            <p>{task.name}</p>
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
            }): Edge[] => {
                return task.spec.dependencies
                    ? task.spec.dependencies.map(
                          (dep: string, i: number): Edge => {
                              return {
                                  id: task.name + i.toString(),
                                  target: dep,
                                  source: task.name,
                              };
                          }
                      )
                    : ([] as Edge[]);
            }
        )
        .flat();

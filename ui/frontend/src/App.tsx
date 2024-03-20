import { PipelineView } from "./PipelineView.tsx";
import { ExecutionView } from "./ExecutionView.tsx";
import { PipelineList } from "./PipelineList.tsx";
import { pipelinesBrambleDev } from "./bramble-types";
import { Pod } from "kubernetes-models/v1";
import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;
import "reactflow/dist/style.css";
import React, { useState, useEffect, useRef } from "react";
import "./index.css";
const NamespaceSearch = (props: {
    inputRef: any;
    setNamespace: any;
}): React.ReactNode => (
    <>
        <div className="label">
            <span className="label-text">Enter Namespace</span>
        </div>
        <input
            className="input input-bordered input-primary"
            type="text"
            ref={props.inputRef}
        />

        <button
            className="btn btn-primary"
            onClick={() => {
                const ns: string | undefined = props.inputRef.current?.value;
                console.log(ns);
                if (ns !== undefined && ns !== "") {
                    props.setNamespace(ns);
                }
            }}
        >
            Get pipelines
        </button>
    </>
);

const App = (): React.ReactNode => {
    const [namespace, setNamespace] = useState<string>("default");

    const [focusedPipeline, setFocusedPipeline] = useState<number>(0);
    const [focusedExecution, setFocusedExecution] = useState<string>("");
    const [podList, setPods] = useState<Pod[]>(new Array<Pod>());

    const [pipelines, setPipelines] = useState<Pipeline[]>(
        new Array<Pipeline>()
    );

    const [executions, setExecutions] = useState<
        pipelinesBrambleDev.v1alpha1.Execution[]
    >(new Array<Execution>());

    const inputRef = useRef<HTMLInputElement>(null);

    const fetchNamespacedResources = async () => {
        const fetchRes = async (resource: string) => {
            console.log("Fetching", resource, "in:", namespace);
            try {
                const baseUrl: string = "http://localhost:5555/";
                await fetch(baseUrl + resource + "/" + namespace)
                    .then((response) => response.json())
                    .then((jsonData) => {
                        switch (resource) {
                            case "pipelines": {
                                setPipelines(jsonData);
                                break;
                            }
                            case "executions": {
                                setExecutions(jsonData);
                                break;
                            }
                            case "pods": {
                                setPods(jsonData);
                                break;
                            }
                        }
                    });
            } catch (error) {
                console.error(error);
            }
        };
        await fetchRes("pipelines");
        await fetchRes("executions");
        await fetchRes("pods");
    };

    useEffect(() => {
        fetchNamespacedResources();
    });

    return (
        <>
            <header className="">
                <h1 className="text-3xl font-bold">Bramble</h1>
            </header>
            <NamespaceSearch inputRef={inputRef} setNamespace={setNamespace} />

            <PipelineList
                namespace={namespace}
                focusedPipeline={focusedPipeline}
                focusedExecution={focusedExecution}
                pipelines={pipelines}
                executions={executions}
                setFocusedPipeline={setFocusedPipeline}
                setFocusedExecution={setFocusedExecution}
            />
            {(!focusedExecution &&
                pipelines.length !== 0 &&
                pipelines[focusedPipeline] && (
                    <PipelineView pipeline={pipelines[focusedPipeline]} />
                )) ||
                (focusedExecution && (
                    <ExecutionView
                        pipeline={pipelines[focusedPipeline]}
                        execution={
                            executions.filter(
                                (exe: Execution) =>
                                    exe.metadata?.name == focusedExecution
                            )[0]
                        }
                        pods={podList}
                    />
                ))}
        </>
    );
};

export default App;

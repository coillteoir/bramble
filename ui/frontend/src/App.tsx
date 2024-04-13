import { ExecutionView } from "./ExecutionView.tsx";
import { PipelineList } from "./PipelineList.tsx";
import { pipelinesBrambleDev } from "./bramble-types";
import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;
import "reactflow/dist/style.css";
import React, { useState, useEffect, useRef } from "react";
import "./index.css";

const NamespaceSearch = (props: {
    inputRef: React.RefObject<HTMLInputElement>;
    fetchNamespacedResources: () => Promise<void>;
    setNamespace: React.Dispatch<React.SetStateAction<string>>;
}): React.ReactNode => (
    <div className="">
        <input
            className="input input-bordered input-primary w-3/5"
            type="text"
            ref={props.inputRef}
            placeholder="Enter Namespace"
        />

        <button
            className="btn btn-primary w-2/5"
            onClick={() => {
                const ns: string | undefined = props.inputRef.current?.value;
                if (ns) {
                    props.setNamespace(ns);
                    return;
                }
                props.fetchNamespacedResources();
            }}
        >
            Get pipelines
        </button>
    </div>
);

const App = (): React.ReactNode => {
    const [namespace, setNamespace] = useState<string>("testns");

    const [focusedPipeline, setFocusedPipeline] = useState<number>(0);
    const [focusedExecution, setFocusedExecution] = useState<
        string | undefined
    >("");

    const [pipelines, setPipelines] = useState<Pipeline[]>(
        new Array<Pipeline>()
    );

    const [executions, setExecutions] = useState<Execution[]>(
        new Array<Execution>()
    );

    const inputRef = useRef<HTMLInputElement>(null);

    const fetchNamespacedResources = async () => {
        const fetchRes = async (resource: "pipelines" | "executions") => {
            console.log("Fetching", resource, "in:", namespace);
            try {
                await fetch(
                    `http://localhost:5555/resources/${resource}/${namespace}`
                )
                    .then((response) => {
                        if (response.status === 200) {
                            return response.json();
                        } else {
                            throw new Error(
                                `error [${response.status}] cannot fetch ${resource} from namespace ${namespace}`
                            );
                        }
                    })
                    .then((jsonData) => {
                        switch (resource) {
                            case "pipelines": {
                                setPipelines(jsonData);
                                console.log("Pipelines set ", jsonData);
                                break;
                            }
                            case "executions": {
                                setExecutions(jsonData);
                                console.log("Executions set ", jsonData);
                                break;
                            }
                        }
                    })
                    .catch((error) => console.error(error));
            } catch (error) {
                console.error(error);
            }
        };
        await fetchRes("pipelines");
        await fetchRes("executions");
    };

    useEffect(() => {
        fetchNamespacedResources();
    }, [namespace]);

    return (
        <>
            <header className="flex">
                <img className="size-10" src="/logo.svg" />
                <h1 className="text-3xl font-bold">Bramble</h1>

            </header>
            <div className="flex h-screen border-4 border-slate-600">
                <div className="inline-block w-1/3 lg:w-1/6">
                    <NamespaceSearch
                        inputRef={inputRef}
                        setNamespace={setNamespace}
                        fetchNamespacedResources={fetchNamespacedResources}
                    />

                    <PipelineList
                        namespace={namespace}
                        pipelines={pipelines}
                        executions={executions}
                        focusedPipeline={focusedPipeline}
                        setFocusedPipeline={setFocusedPipeline}
                        setFocusedExecution={setFocusedExecution}
                    />
                </div>
                <div className="inline-block w-2/3 border-2 border-white lg:w-5/6">
                    {pipelines[focusedPipeline] && (
                        <ExecutionView
                            pipeline={pipelines[focusedPipeline]}
                            execution={executions.find(
                                (exe: Execution) =>
                                    exe.metadata?.name === focusedExecution
                            )}
                        />
                    )}
                </div>
            </div>
        </>
    );
};

export default App;

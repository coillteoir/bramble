import { pipelinesBrambleDev } from "./bramble-types";
import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;
import React from "react";
import "./index.css";

const ExecutionList = (props: {
    namespace: string;
    focusedPipeline: number;
    focusedExecution: string | undefined;
    pipeline: Pipeline;
    pipelines: Array<Pipeline>;
    executions: Array<Execution>;
    setFocusedPipeline: React.Dispatch<React.SetStateAction<number>>;
    setFocusedExecution: React.Dispatch<
        React.SetStateAction<string | undefined>
    >;
}): React.ReactNode => (
    <>
        <p className="">{props.pipeline.metadata?.name}</p>
        {props.pipeline === props.pipelines[props.focusedPipeline] && (
            <ul className="bg-slate-600">
                {props.executions
                    .filter(
                        (exe: Execution) =>
                            exe.spec?.pipeline === props.pipeline.metadata?.name
                    )
                    .map((exe: Execution, i: number) => (
                        <li
                            onClick={() =>
                                props.setFocusedExecution(exe.metadata?.name)
                            }
                            key={i}
                        >
                            {exe.metadata?.name}
                        </li>
                    ))}
            </ul>
        )}
    </>
);

export const PipelineList = (props: {
    namespace: string;
    focusedPipeline: number;
    focusedExecution: string | undefined;
    pipelines: Array<Pipeline>;
    executions: Array<Execution>;
    setFocusedPipeline: React.Dispatch<React.SetStateAction<number>>;
    setFocusedExecution: React.Dispatch<
        React.SetStateAction<string | undefined>
    >;
}): React.ReactNode => (
    <div className="">
        {props.pipelines.length !== 0 && (
            <h1>Pipelines in the {props.namespace} namespace</h1>
        )}
        <ul className="">
            {props.pipelines.map(
                (pipeline: Pipeline, index: number) =>
                    pipeline && (
                        <li
                            key={index}
                            className=""
                            onClick={() => {
                                props.setFocusedPipeline(index);
                                console.log(
                                    props.pipelines[props.focusedPipeline]
                                );
                            }}
                        >
                            <ExecutionList
                                namespace={props.namespace}
                                focusedPipeline={props.focusedPipeline}
                                focusedExecution={props.focusedExecution}
                                pipeline={pipeline}
                                pipelines={props.pipelines}
                                executions={props.executions}
                                setFocusedPipeline={props.setFocusedPipeline}
                                setFocusedExecution={props.setFocusedExecution}
                            />
                        </li>
                    )
            )}
        </ul>
    </div>
);

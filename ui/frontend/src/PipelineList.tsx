import { pipelinesBrambleDev } from "./bramble-types";
import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;
import React from "react";
import "./index.css";

export const PipelineList = (props: {
    namespace: string;
    focusedPipeline: number;
    focusedExecution: string;
    pipelines: Array<Pipeline>;
    executions: Array<Execution>;
    setFocusedPipeline: any;
    setFocusedExecution: any;
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
                            <p className="">{pipeline.metadata?.name}</p>
                            {pipeline ==
                                props.pipelines[props.focusedPipeline] && (
                                <ul className="bg-slate-600">
                                    {props.executions
                                        .filter(
                                            (exe: Execution) =>
                                                exe.spec?.pipeline ==
                                                pipeline.metadata?.name
                                        )
                                        .map((exe: Execution) => (
                                            <li
                                                onClick={() =>
                                                    props.setFocusedExecution(
                                                        exe.metadata?.name
                                                    )
                                                }
                                            >
                                                {exe.metadata?.name}
                                            </li>
                                        ))}
                                </ul>
                            )}
                        </li>
                    )
            )}
        </ul>
    </div>
);

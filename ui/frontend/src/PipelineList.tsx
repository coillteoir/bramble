import { pipelinesBrambleDev } from "./bramble-types";

import Pipeline = pipelinesBrambleDev.v1alpha1.Pipeline;
import Execution = pipelinesBrambleDev.v1alpha1.Execution;
import React from "react";
import "./index.css";

import { ExecutionPhase, PipelineStatusIcon } from "./ExecutionGraph";

const ExecutionStatusIcon = (props: { phase: string | undefined }) => {
    switch (props.phase) {
        case "running":
            return <PipelineStatusIcon phase={ExecutionPhase.Running} />;
        case "error":
            return <PipelineStatusIcon phase={ExecutionPhase.Failure} />;
        case "completed":
            return <PipelineStatusIcon phase={ExecutionPhase.Success} />;
        default:
            return <PipelineStatusIcon phase={ExecutionPhase.Pending} />;
    }
};

const ExecutionList = (props: {
    pipeline: string | undefined;
    executions: Array<Execution>;
    setFocusedExecution: React.Dispatch<
        React.SetStateAction<string | undefined>
    >;
}): React.ReactNode => (
    <ul className="menu rounded-box bg-slate-700">
        {props.executions
            .filter((exe: Execution) => exe.spec?.pipeline === props.pipeline)
            .map((exe: Execution, i: number) => (
                <li
                    onClick={() =>
                        props.setFocusedExecution(exe.metadata?.name)
                    }
                    key={i}
                >
                    <div className="flex">
                        <p className="inline-block">{exe.metadata?.name}</p>
                        <ExecutionStatusIcon phase={exe.status?.phase} />
                    </div>
                </li>
            ))}
    </ul>
);

export const PipelineList = (props: {
    namespace: string;
    focusedPipeline: number;
    pipelines: Array<Pipeline>;
    executions: Array<Execution>;
    setFocusedPipeline: React.Dispatch<React.SetStateAction<number>>;
    setFocusedExecution: React.Dispatch<
        React.SetStateAction<string | undefined>
    >;
}): React.ReactNode => (
    <div className="inline-block">
        {props.pipelines.length !== 0 && (
            <>
                <h1>Pipelines in the {props.namespace} namespace</h1>

                <ul className="menu menu-vertical rounded-box">
                    {props.pipelines.map(
                        (pipeline: Pipeline, index: number) => (
                            <li
                                key={index}
                                className=""
                                onClick={() => {
                                    props.setFocusedPipeline(index);
                                }}
                            >
                                <p
                                    className=""
                                    onClick={() => {
                                        props.setFocusedExecution("");
                                    }}
                                >
                                    {pipeline.metadata?.name}
                                </p>
                                {index === props.focusedPipeline && (
                                    <ExecutionList
                                        pipeline={pipeline.metadata?.name}
                                        executions={props.executions}
                                        setFocusedExecution={
                                            props.setFocusedExecution
                                        }
                                    />
                                )}
                            </li>
                        )
                    )}
                </ul>
            </>
        )}
    </div>
);

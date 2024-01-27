import { pipelinesBrambleDev } from "./bramble-types";
import { Component, For } from "solid-js";

const PipelineView: Component<{
    pipeline: pipelinesBrambleDev.v1alpha1.Pipeline;
}> = (props: { pipeline: pipelinesBrambleDev.v1alpha1.Pipeline }) => {
    const pl: pipelinesBrambleDev.v1alpha1.Pipeline = props.pipeline;
    console.log(pl);
    return (
        <div class="">
            <h2 class="">Pipeline: {pl.metadata?.name}</h2>
            <h2 class=""> Tasks </h2>
            <ul class="">
                {pl?.spec?.tasks && (
                    <PLTaskView task={pl?.spec?.tasks[0]} pipeline={pl} />
                )}
            </ul>
        </div>
    );
};

const PLTaskView: Component<{
    task: any;
    pipeline: pipelinesBrambleDev.v1alpha1.Pipeline;
}> = (props: {
    task: any;
    pipeline: pipelinesBrambleDev.v1alpha1.Pipeline;
}) => {
    const task: any = props.task;
    const pipeline: pipelinesBrambleDev.v1alpha1.Pipeline = props.pipeline;
    return (
        <div class="">
            {!task.spec.dependencies && <h3 class="">{task.name}</h3>}
            {task.spec.dependencies && (
                <details class="">
                    <summary>{task.name}</summary>
                    <ul>
                        <For each={task.spec.dependencies}>
                            {(dep: string) => {
                                const task0 = pipeline?.spec?.tasks?.filter(
                                    (task) => task.name == dep
                                )[0];

                                return (
                                    task0 && (
                                        <li class="">
                                            <PLTaskView
                                                task={task0}
                                                pipeline={pipeline}
                                            />
                                        </li>
                                    )
                                );
                            }}
                        </For>
                    </ul>
                </details>
            )}
        </div>
    );
};

export { PipelineView };

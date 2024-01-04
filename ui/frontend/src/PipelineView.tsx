import { Pipeline, PLtask } from "./bramble_types.ts";
import { Component, For } from "solid-js";

/*
 * TODO: Redefine Tasks to recursively generate tree diagram.
 */

const PipelineView: Component<{ pipeline: Pipeline }> = (props: {
    pipeline: Pipeline;
}) => {
    const pl: Pipeline = props.pipeline;
    console.log(pl);
    return (
        <div class="grid gap-4 pipeline rounded gap">
            <h2 class="text-2xl font-bold">Name: {pl.metadata.name}</h2>
            <h2 class="text-2xl font-bold"> Tasks </h2>
            <ul class="bg-gray grid gap-3">
                {pl.spec.tasks && (
                    <PLTaskView task={pl.spec?.tasks[0]} pipeline={pl} />
                )}
            </ul>
        </div>
    );
};

const PLTaskView: Component<{ task: PLtask; pipeline: Pipeline }> = (props: {
    task: PLtask;
    pipeline: Pipeline;
}) => {
    const task: PLtask = props.task;
    const pipeline: Pipeline = props.pipeline;
    return (
        <div
            class="grid rounded bg-slate-500 
        border-2 border-white border-solid"
        >
            <h3 class="bold text-lg font-medium">{task.name}</h3>
            {task.spec.dependencies && (
                <>
                    <ul>
                        <For each={task.spec.dependencies}>
                            {(dep: string) => {
                                const task0 = pipeline.spec.tasks?.filter(
                                    (task) => task.name == dep
                                )[0];

                                return (
                                    task0 && (
                                        <PLTaskView
                                            task={task0}
                                            pipeline={pipeline}
                                        />
                                    )
                                );
                            }}
                        </For>
                    </ul>
                </>
            )}
        </div>
    );
};

export { PipelineView };

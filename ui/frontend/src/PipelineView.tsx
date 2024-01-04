import { Pipeline, PLtask } from "./bramble_types.ts";
import { Component, For } from "solid-js";

const PipelineView: Component<{ pipeline: Pipeline }> = (props: {
    pipeline: Pipeline;
}) => {
    const pl: Pipeline = props.pipeline;
    console.log(pl);
    return (
        <div class="grid gap-4 pipeline rounded bg-black gap">
            <h2 class="text-2xl font-bold">Name: {pl?.metadata.name}</h2>
            <h2 class="text-2xl font-bold"> Tasks </h2>
            <ul class="bg-gray grid gap-3">
                <For each={pl.spec.tasks}>
                    {(task: PLtask) => task && <PLTaskView task={task} />}
                </For>
            </ul>
        </div>
    );
};

const PLTaskView: Component<{ task: PLtask }> = (props: { task: PLtask }) => {
    const task: PLtask = props?.task;
    return (
        <div class="grid task bg-gray-600">
            <h3 class="bold text-lg font-medium">{task?.name}</h3>
            {task.spec.dependencies && (
                <>
                    <h3>Dependencies</h3>
                    <ul>
                        <For each={task.spec.dependencies}>
                            {(dep: string) => <li>{dep}</li>}
                        </For>
                    </ul>
                </>
            )}
        </div>
    );
};

export { PipelineView };

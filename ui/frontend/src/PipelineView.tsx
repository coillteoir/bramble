import { Pipeline, PLtask } from "./bramble_types.ts";
import { Component, For } from "solid-js";

const PipelineView: Component<{ pipeline: Pipeline }> = (props: {
    pipeline: Pipeline;
}) => {
    const pl: Pipeline = props.pipeline;
    console.log(pl);
    return (
        <div class="p-8 bg-black border-2 border-white">
            <h2 class="">Pipeline: {pl.metadata.name}</h2>
            <h2 class=""> Tasks </h2>
            <ul class="">
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
        <div class="rounded-lg w-auto border-b-2">
            {!task.spec.dependencies && <h3 class="">{task.name}</h3>}
            {task.spec.dependencies && (
                <details open class="">
                    <summary>{task.name}</summary>
                    <ul>
                        <For each={task.spec.dependencies}>
                            {(dep: string) => {
                                const task0 = pipeline.spec.tasks?.filter(
                                    (task) => task.name == dep
                                )[0];

                                return (
                                    task0 && (
                                        <li class="border-l-2 border-l-white pl-8">
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

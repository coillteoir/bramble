import { Pipeline, PLtask } from "./bramble_types.ts";
import { Component } from "solid-js";

const PipelineView: Component<{ pipeline: Pipeline}> = (props: {
  pipeline: Pipeline;
}) => {
  const pl: Pipeline = props?.pipeline;
  console.log(pl);
  return (
    <div class="grid gap-4 pipeline rounded bg-black gap">
      <h2 class="text-2xl font-bold">Name: {pl?.metadata.name}</h2>
      <h2 class="text-2xl font-bold"> Tasks </h2>
      <ul class="bg-gray grid gap-3">
        {pl?.spec.tasks?.map((task: PLtask) => <PLTaskView task={task} />)}
      </ul>
    </div>
  );
};

const PLTaskView: Component<{ task: PLtask | undefined }> = (props: {
  task: PLtask | undefined;
}) => {
  const task: PLtask | undefined = props?.task;
  return (
    <div class="grid task bg-gray-600">
      <h3 class="bold text-lg font-medium">{task?.name}</h3>
      {task?.spec.dependencies && (
        <>
          <h3>Dependencies</h3>
          <ul>
            {task?.spec.dependencies?.map((dep: string) => <li>{dep}</li>)}
          </ul>
        </>
      )}
    </div>
  );
};

export { PipelineView };

import { Pipeline, PLtask } from "./bramble_types.ts";
import { Component } from "solid-js";

const PipelineView: Component<{pipeline: Pipeline | undefined}> = (props: {pipeline: Pipeline | undefined}) => {
  console.log(props.pipeline)
  const pl: Pipeline | undefined = props?.pipeline
  return (
      <div class="pipeline">
        <h2>Name: {pl?.metadata.name}</h2>
        <h2> Tasks </h2>
        <ul>
        {pl?.spec.tasks?.map((task: PLtask) => <PLTaskView task={task} />)}
        </ul>
      </div>
  );
};

const PLTaskView: Component<{task: PLtask | undefined}> = (props: {task: PLtask | undefined}) => {
  const task :PLtask | undefined = props?.task
  return (
      
      <div class="task">
      <h3>{task?.name}</h3>
      {task?.dependencies && (
          <>
            <h3>Dependencies</h3>
            <ul>
                {task?.dependencies?.map((dep: string) => <li>{dep}</li>)}
            </ul>
        </>
      )}
          </div>
  );
};

export { PipelineView };

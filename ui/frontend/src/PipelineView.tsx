import { Pipeline, PLtask } from "./bramble_types.ts";
import {Component} from "solid-js"

const PipelineView : Component<{pipeline: Pipeline}> = (props) => {
    return (
    <div>
        {props}
        <PLTaskView />
    </div>
    )
}

const PLTaskView : Component<{task: PLtask}> = (props) => {
    return (
       <div>{props}</div> 
       )
}

export {PipelineView}

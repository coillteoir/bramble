"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Pipeline = exports.TaskRef = exports.PLtask = exports.TaskSpec = void 0;
var TaskSpec = /** @class */ (function () {
    function TaskSpec(image, command) {
        this.image = image;
        this.command = command;
    }
    return TaskSpec;
}());
exports.TaskSpec = TaskSpec;
var PLtask = /** @class */ (function () {
    function PLtask(name, spec, dependencies) {
        this.name = name;
        this.spec = spec;
        this.dependencies = dependencies;
    }
    return PLtask;
}());
exports.PLtask = PLtask;
var TaskRef = /** @class */ (function () {
    function TaskRef(name, dependencies) {
        this.name = name;
        this.dependencies = dependencies;
    }
    return TaskRef;
}());
exports.TaskRef = TaskRef;
var Pipeline = /** @class */ (function () {
    function Pipeline(metadata, spec) {
        this.metadata = metadata;
        this.spec = spec;
    }
    return Pipeline;
}());
exports.Pipeline = Pipeline;

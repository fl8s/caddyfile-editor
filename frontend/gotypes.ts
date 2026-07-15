// AUTOMATICALLY GENERATED - DO NOT EDIT

export interface UpstreamStatus {
	"address": string;
	"up": boolean;
	"error"?: string;
}
export interface LivenessStatus {
	"host": string;
	"upstreams": UpstreamStatus[];
}

export interface Warning {
	"file"?: string;
	"line"?: number /* int */;
	"directive"?: string;
	"message"?: string;
}
export interface AdaptResult {
	"Warnings": Warning[];
	"AdaptError"?: string;
}

/// RPC Generated
export interface App {
	AdaptCaddyfile(arg0: string): Promise<AdaptResult>
	InstallCaddyfile(arg0: string): Promise<boolean>
	LastCaddyfile(): Promise<string>
    Liveness(): Promise<LivenessStatus[]>
}

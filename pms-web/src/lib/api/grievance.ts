import { get, post } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type { Grievance } from "@/types/performance";

export const getStaffGrievances = (staffId: string) =>
  get<BaseAPIResponse<Grievance[]>>(`/grievances?staffId=${staffId}`);

export const getGrievanceReport = () =>
  get<BaseAPIResponse<Grievance[]>>("/grievances/report");

export const submitGrievance = (data: unknown) =>
  post<ResponseVm>("/grievances", data);

export const submitGrievanceResolution = (data: unknown) =>
  post<ResponseVm>("/grievances/resolution", data);

export const escalateGrievance = (data: unknown) =>
  post<ResponseVm>("/grievances/escalate", data);

export const closeGrievance = (data: unknown) =>
  post<ResponseVm>("/grievances/close", data);

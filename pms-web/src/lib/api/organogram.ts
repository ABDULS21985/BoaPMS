import { get, post } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type { Directorate, Department, Division, Office } from "@/types/organogram";

export const getDirectorates = () =>
  get<BaseAPIResponse<Directorate[]>>("/organogram/directorates");
export const saveDirectorate = (data: unknown) =>
  post<ResponseVm>("/organogram/directorates", data);

export const getDepartments = () =>
  get<BaseAPIResponse<Department[]>>("/organogram/departments");
export const saveDepartment = (data: unknown) =>
  post<ResponseVm>("/organogram/departments", data);

export const getDivisions = () =>
  get<BaseAPIResponse<Division[]>>("/organogram/divisions");
export const getDivisionsByDepartment = (departmentId: number) =>
  get<BaseAPIResponse<Division[]>>(`/organogram/departments/${departmentId}/divisions`);
export const saveDivision = (data: unknown) =>
  post<ResponseVm>("/organogram/divisions", data);

export const getOffices = () =>
  get<BaseAPIResponse<Office[]>>("/organogram/offices");
export const getOfficesByDivision = (divisionId: number) =>
  get<BaseAPIResponse<Office[]>>(`/organogram/divisions/${divisionId}/offices`);
export const saveOffice = (data: unknown) =>
  post<ResponseVm>("/organogram/offices", data);

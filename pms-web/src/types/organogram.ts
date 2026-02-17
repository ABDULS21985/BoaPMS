export interface Directorate {
  directorateId: number;
  directorateName: string;
  directorateCode?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Department {
  departmentId: number;
  directorateId?: number;
  departmentName: string;
  departmentCode?: string;
  isBranch: boolean;
  directorateName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Division {
  divisionId: number;
  departmentId: number;
  divisionName: string;
  divisionCode?: string;
  departmentName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Office {
  officeId: number;
  divisionId: number;
  officeName: string;
  officeCode?: string;
  divisionName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface BankYear {
  bankYearId: number;
  year: number;
  name: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

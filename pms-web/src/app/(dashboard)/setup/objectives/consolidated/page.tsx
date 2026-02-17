"use client";

import { useEffect, useState, useCallback } from "react";
import { Search, ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { PageHeader } from "@/components/shared/page-header";
import { PageSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { StatusBadge } from "@/components/shared/status-badge";
import { getConsolidatedObjectivesPaginated } from "@/lib/api/performance";
import { getDepartments, getDivisionsByDepartment, getOfficesByDivision } from "@/lib/api/organogram";
import type { ConsolidatedObjective } from "@/types/performance";
import type { Department, Division, Office } from "@/types/organogram";

const levelColors: Record<string, string> = {
  Enterprise: "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300",
  Department: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
  Division: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
  Office: "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300",
};

export default function ConsolidatedObjectivesPage() {
  const [items, setItems] = useState<ConsolidatedObjective[]>([]);
  const [totalRecords, setTotalRecords] = useState(0);
  const [loading, setLoading] = useState(true);
  const [pageIndex, setPageIndex] = useState(0);
  const [pageSize, setPageSize] = useState(20);
  const [search, setSearch] = useState("");
  const [departments, setDepartments] = useState<Department[]>([]);
  const [divisions, setDivisions] = useState<Division[]>([]);
  const [offices, setOffices] = useState<Office[]>([]);
  const [filterDept, setFilterDept] = useState("");
  const [filterDiv, setFilterDiv] = useState("");
  const [filterOffice, setFilterOffice] = useState("");

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const res = await getConsolidatedObjectivesPaginated({
        pageIndex,
        pageSize,
        searchString: search || undefined,
        departmentId: filterDept ? Number(filterDept) : undefined,
        divisionId: filterDiv ? Number(filterDiv) : undefined,
        officeId: filterOffice ? Number(filterOffice) : undefined,
      });
      if (res?.data) {
        setItems(res.data.objectives ?? []);
        setTotalRecords(res.data.totalRecords ?? 0);
      }
    } catch { /* */ } finally { setLoading(false); }
  }, [pageIndex, pageSize, search, filterDept, filterDiv, filterOffice]);

  useEffect(() => { loadData(); }, [loadData]);

  useEffect(() => {
    getDepartments().then((r) => { if (r?.data) setDepartments(Array.isArray(r.data) ? r.data : []); });
  }, []);

  const onDeptFilter = async (val: string) => {
    setFilterDept(val); setFilterDiv(""); setFilterOffice(""); setDivisions([]); setOffices([]); setPageIndex(0);
    if (val) { try { const r = await getDivisionsByDepartment(Number(val)); if (r?.data) setDivisions(Array.isArray(r.data) ? r.data : []); } catch { /* */ } }
  };

  const onDivFilter = async (val: string) => {
    setFilterDiv(val); setFilterOffice(""); setOffices([]); setPageIndex(0);
    if (val) { try { const r = await getOfficesByDivision(Number(val)); if (r?.data) setOffices(Array.isArray(r.data) ? r.data : []); } catch { /* */ } }
  };

  const clearFilters = () => {
    setFilterDept(""); setFilterDiv(""); setFilterOffice(""); setDivisions([]); setOffices([]); setSearch(""); setPageIndex(0);
  };

  const totalPages = Math.max(1, Math.ceil(totalRecords / pageSize));

  const headerCols = ["Objective Name", "Level", "SBU", "KPI", "Target", "Status"];

  if (loading && items.length === 0) return <div><PageHeader title="Consolidated Objectives" breadcrumbs={[{ label: "Setup" }, { label: "Consolidated Objectives" }]} /><PageSkeleton /></div>;

  return (
    <div className="space-y-6">
      <PageHeader title="Consolidated Objectives" description="View all objectives across enterprise, department, division, and office levels" breadcrumbs={[{ label: "Setup", href: "/setup/strategies" }, { label: "Consolidated Objectives" }]} />

      {/* Filters */}
      <div className="flex flex-wrap items-end gap-3">
        <div className="relative w-64">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input placeholder="Search objectives..." value={search} onChange={(e) => { setSearch(e.target.value); setPageIndex(0); }} className="pl-9" />
        </div>
        <Select value={filterDept} onValueChange={onDeptFilter}>
          <SelectTrigger className="w-48"><SelectValue placeholder="All Departments" /></SelectTrigger>
          <SelectContent>{departments.map((d) => <SelectItem key={d.departmentId} value={String(d.departmentId)}>{d.departmentName}</SelectItem>)}</SelectContent>
        </Select>
        {divisions.length > 0 && (
          <Select value={filterDiv} onValueChange={onDivFilter}>
            <SelectTrigger className="w-48"><SelectValue placeholder="All Divisions" /></SelectTrigger>
            <SelectContent>{divisions.map((d) => <SelectItem key={d.divisionId} value={String(d.divisionId)}>{d.divisionName}</SelectItem>)}</SelectContent>
          </Select>
        )}
        {offices.length > 0 && (
          <Select value={filterOffice} onValueChange={(v) => { setFilterOffice(v); setPageIndex(0); }}>
            <SelectTrigger className="w-48"><SelectValue placeholder="All Offices" /></SelectTrigger>
            <SelectContent>{offices.map((o) => <SelectItem key={o.officeId} value={String(o.officeId)}>{o.officeName}</SelectItem>)}</SelectContent>
          </Select>
        )}
        {(filterDept || filterDiv || filterOffice || search) && (
          <Button variant="ghost" size="sm" onClick={clearFilters}>Clear</Button>
        )}
      </div>

      {/* Table */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              {headerCols.map((h) => <TableHead key={h}>{h}</TableHead>)}
            </TableRow>
          </TableHeader>
          <TableBody>
            {items.length > 0 ? items.map((item) => (
              <TableRow key={item.objectiveId}>
                <TableCell className="font-medium">{item.name}</TableCell>
                <TableCell><Badge className={levelColors[item.objectiveLevel] ?? ""}>{item.objectiveLevel}</Badge></TableCell>
                <TableCell>{item.sbuName ?? "—"}</TableCell>
                <TableCell><span className="line-clamp-1">{item.kpi ?? "—"}</span></TableCell>
                <TableCell>{item.target ?? "—"}</TableCell>
                <TableCell>{item.recordStatus != null ? <StatusBadge status={item.recordStatus} /> : "—"}</TableCell>
              </TableRow>
            )) : (
              <TableRow><TableCell colSpan={headerCols.length} className="h-24 text-center"><EmptyState title="No objectives found" /></TableCell></TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">{totalRecords} objective(s) total</p>
        <div className="flex items-center gap-2">
          <Select value={String(pageSize)} onValueChange={(v) => { setPageSize(Number(v)); setPageIndex(0); }}>
            <SelectTrigger className="h-8 w-[70px]"><SelectValue /></SelectTrigger>
            <SelectContent>{[10, 20, 50, 100].map((s) => <SelectItem key={s} value={String(s)}>{s}</SelectItem>)}</SelectContent>
          </Select>
          <span className="text-sm text-muted-foreground">Page {pageIndex + 1} of {totalPages}</span>
          <Button variant="outline" size="icon" className="h-8 w-8" onClick={() => setPageIndex((p) => p - 1)} disabled={pageIndex === 0}><ChevronLeft className="h-4 w-4" /></Button>
          <Button variant="outline" size="icon" className="h-8 w-8" onClick={() => setPageIndex((p) => p + 1)} disabled={pageIndex + 1 >= totalPages}><ChevronRight className="h-4 w-4" /></Button>
        </div>
      </div>
    </div>
  );
}

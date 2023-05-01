import { SmallSpan } from "../globals";
import Placeholder from "../placeholder";
import React from "react";
import {
  Column,
  Row,
  useGlobalFilter,
  usePagination,
  useTable,
} from "react-table";
import {
  StyledTd,
  StyledTable,
  StyledTHead,
  StyledTh,
  StyledTBody,
} from "./styles";

export type TableProps = {
  columns: Column<any>[];
  data: any[];
  dataName?: string;
  onRowClick?: (row: Row) => void;
  rowHeight?: string;
};

const Table: React.FC<TableProps> = ({
  columns: columnsData,
  data,
  dataName,
  onRowClick,
  rowHeight,
}) => {
  const { rows, getTableProps, getTableBodyProps, prepareRow, headerGroups } =
    useTable(
      {
        columns: columnsData,
        data,
      },
      useGlobalFilter,
      usePagination
    );

  const renderRows = () => {
    return (
      <>
        {rows.map((row: any) => {
          prepareRow(row);

          return (
            <tr
              {...row.getRowProps()}
              enablepointer={(!!onRowClick).toString()}
              onClick={() => onRowClick && onRowClick(row)}
              selected={false}
            >
              {row.cells.map((cell: any) => {
                return (
                  <StyledTd
                    {...cell.getCellProps()}
                    style={{
                      width: cell.column.totalWidth,
                    }}
                  >
                    {cell.render("Cell")}
                  </StyledTd>
                );
              })}
            </tr>
          );
        })}
      </>
    );
  };

  if (rows.length == 0) {
    return (
      <Placeholder>
        <SmallSpan>No {dataName || "data"} found.</SmallSpan>
      </Placeholder>
    );
  }

  return (
    <>
      <StyledTable {...getTableProps()}>
        <StyledTHead>
          {headerGroups.map((headerGroup) => (
            <tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <StyledTh {...column.getHeaderProps()}>
                  {column.render("Header")}
                </StyledTh>
              ))}
            </tr>
          ))}
        </StyledTHead>
        <StyledTBody
          rowHeight={rowHeight}
          {...getTableBodyProps()}
          canClick={!!onRowClick}
        >
          {renderRows()}
        </StyledTBody>
      </StyledTable>
    </>
  );
};

export default Table;

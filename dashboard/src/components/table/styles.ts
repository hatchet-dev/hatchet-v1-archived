import { tableHeaderFontStack, textFontStack } from "components/globals";
import theme from "shared/theme";
import styled from "styled-components";

export const StyledTable = styled.table`
  width: 100%;
  min-width: 500px;
  border-spacing: 0;
`;

export const StyledTHead = styled.thead`
  width: 100%;
  position: sticky;

  > tr {
    background: ${theme.bg.shadeone};
    line-height: 2.5em;

    > th {
      ${tableHeaderFontStack};
      border-top: ${theme.line.default};
      border-bottom: ${theme.line.default};
    }
  }

  > tr:first-child {
    > th:first-child {
      border-top-left-radius: 6px;
      border-left: ${theme.line.default};
    }

    > th:last-child {
      border-top-right-radius: 6px;
      border-right: ${theme.line.default};
    }
  }
`;

export const StyledTBody = styled.tbody<{ rowHeight?: string }>`
  > tr {
    background: ${theme.bg.shadeone};
    line-height: ${(props) => props.rowHeight || "2.5em"};
    cursor: pointer;
    :hover {
      background: ${theme.bg.hover};
    }

    > td {
      border-bottom: ${theme.line.default};
    }

    > td:first-child {
      border-left: ${theme.line.default};
    }

    > td:last-child {
      border-right: ${theme.line.default};
    }
  }

  > tr:last-child {
    > td:first-child {
      border-bottom-left-radius: 6px;
      border-left: ${theme.line.default};
    }

    > td:last-child {
      border-bottom-right-radius: 6px;
      border-right: ${theme.line.default};
    }
  }
`;

export const StyledTd = styled.td`
  ${textFontStack}
  font-size: 13px;
  color: ${theme.text.default};
  :first-child {
    padding-left: 20px;
  }

  :last-child {
    padding-right: 20px;
  }

  user-select: text;
`;

export const StyledTh = styled.th`
  ${textFontStack}

  text-align: left;
  font-size: 13px;
  font-weight: 500;
  color: ${theme.text.default};
  :first-child {
    padding-left: 20px;
  }
  :last-child {
    padding-right: 20px;
  }
`;

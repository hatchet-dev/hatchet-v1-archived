import { tableHeaderFontStack, textFontStack } from "components/globals";
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
    background: ${(props) => props.theme.bg.shadeone};
    line-height: 2.5em;

    > th {
      ${tableHeaderFontStack};
      border-top: ${(props) => props.theme.line.default};
      border-bottom: ${(props) => props.theme.line.default};
    }
  }

  > tr:first-child {
    > th:first-child {
      border-top-left-radius: 6px;
      border-left: ${(props) => props.theme.line.default};
    }

    > th:last-child {
      border-top-right-radius: 6px;
      border-right: ${(props) => props.theme.line.default};
    }
  }
`;

export const StyledTBody = styled.tbody<{
  rowHeight?: string;
  canClick: boolean;
}>`
  > tr {
    background: ${(props) => props.theme.bg.shadeone};
    line-height: ${(props) => props.rowHeight || "2.5em"};
    cursor: ${(props) => (props.canClick ? "pointer" : "default")};

    ${(props) =>
      props.canClick &&
      `
    :hover {
      background: ${props.theme.bg.hover};
    }
    `}

    > td {
      border-bottom: ${(props) => props.theme.line.default};
    }

    > td:first-child {
      border-left: ${(props) => props.theme.line.default};
    }

    > td:last-child {
      border-right: ${(props) => props.theme.line.default};
    }
  }

  > tr:last-child {
    > td:first-child {
      border-bottom-left-radius: 6px;
      border-left: ${(props) => props.theme.line.default};
    }

    > td:last-child {
      border-bottom-right-radius: 6px;
      border-right: ${(props) => props.theme.line.default};
    }
  }
`;

export const StyledTd = styled.td`
  ${textFontStack}
  font-size: 13px;
  color: ${(props) => props.theme.text.default};
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
  color: ${(props) => props.theme.text.default};
  :first-child {
    padding-left: 20px;
  }
  :last-child {
    padding-right: 20px;
  }
`;

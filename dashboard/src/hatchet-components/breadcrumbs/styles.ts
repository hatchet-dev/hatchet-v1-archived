import { FlexRow, MaterialIcon } from "../globals";
import styled from "styled-components";

export const BreadcrumbWrapper = styled(FlexRow)`
  max-width: fit-content;
  margin: 14px 0;
`;

export const Breadcrumb = styled.a<{
  clickable: boolean;
}>`
  color: ${(props) => props.theme.text.default};
  font-weight: ${(props) => (props.clickable ? "700" : "500")};
  font-size: 11px;
  cursor: ${(props) => (props.clickable ? "pointer" : "default")};
`;

export const BreadcrumbArrow = styled(MaterialIcon)`
  color: ${(props) => props.theme.text.default};
  margin: 0 12px;
  font-size: 14px;
`;

import { FlexRow, tableHeaderFontStack } from "hatchet-components";
import styled from "styled-components";

export const InviteContainer = styled(FlexRow)`
  gap: 12px;
`;

export const InviteAddButton = styled(FlexRow)<{
  size?: "small" | "default";
  has_icon?: boolean;
  disabled?: boolean;
}>`
  ${tableHeaderFontStack};
  padding: 0 6px;
  background-color: ${(props) =>
    props.disabled ? props.theme.bg.inactive : props.theme.bg.buttonone};
  color: white;
  border-radius: 50%;
  cursor: ${(props) => (props.disabled ? "default" : "pointer")};

  > i {
    font-weight: 700;
  }

  :hover {
    filter: brightness(93%);
  }
`;

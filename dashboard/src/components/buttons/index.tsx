import { MaterialIcon } from "components/globals";
import React from "react";
import { StyledStandardButton } from "./styles";

export type Props = {
  material_icon?: string;
  label: string;
  size?: "small" | "default";
  icon_side?: "left" | "right";
  margin?: string;
  on_click?: () => void;
};

const StandardButton: React.FC<Props> = ({
  material_icon,
  label,
  size = "default",
  icon_side = "left",
  margin = "0 8px",
  on_click,
}) => {
  const children = [
    <MaterialIcon key="0" className="material-icons">
      {material_icon}
    </MaterialIcon>,
    label,
  ];

  if (material_icon) {
    return (
      <StyledStandardButton
        size={size}
        has_icon={true}
        onClick={on_click}
        icon_side={icon_side}
        margin={margin}
      >
        {icon_side == "left" ? children : children.reverse()}
      </StyledStandardButton>
    );
  }

  return (
    <StyledStandardButton size={size} has_icon={false} onClick={on_click}>
      {label}
    </StyledStandardButton>
  );
};

export default StandardButton;

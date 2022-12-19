import { MaterialIcon } from "components/globals";
import React from "react";
import { StyledStandardButton } from "./styles";

export type Props = {
  material_icon?: string;
  label: string;
  size?: "small" | "default";
  icon_side?: "left" | "right";
  on_click?: () => void;
};

const StandardButton: React.FC<Props> = ({
  material_icon,
  label,
  size = "default",
  icon_side = "left",
  on_click,
}) => {
  const children = [
    <MaterialIcon className="material-icons">{material_icon}</MaterialIcon>,
    label,
  ];

  if (material_icon) {
    return (
      <StyledStandardButton
        size={size}
        has_icon={true}
        onClick={on_click}
        icon_side={icon_side}
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

import { MaterialIcon } from "components/globals";
import Spinner from "components/loaders";
import React from "react";
import { StyledStandardButton } from "./styles";

export type Props = {
  material_icon?: string;
  label: string;
  size?: "small" | "default";
  icon_side?: "left" | "right";
  margin?: string;
  disabled?: boolean;
  on_click?: () => void;
  is_loading?: boolean;
};

const StandardButton: React.FC<Props> = ({
  material_icon,
  label,
  size = "default",
  icon_side = "left",
  margin = "0 8px",
  disabled = false,
  on_click,
  is_loading = false,
}) => {
  if (is_loading) {
    const children = [<Spinner />, label];

    return (
      <StyledStandardButton
        size={size}
        has_icon={true}
        onClick={on_click}
        icon_side={icon_side}
        margin={margin}
        disabled={disabled}
      >
        {icon_side == "left" ? children : children.reverse()}
      </StyledStandardButton>
    );
  }

  if (material_icon) {
    const children = [
      <MaterialIcon key="0" className="material-icons">
        {material_icon}
      </MaterialIcon>,
      label,
    ];

    return (
      <StyledStandardButton
        size={size}
        has_icon={true}
        onClick={on_click}
        icon_side={icon_side}
        margin={margin}
        disabled={disabled}
      >
        {icon_side == "left" ? children : children.reverse()}
      </StyledStandardButton>
    );
  }

  return (
    <StyledStandardButton
      size={size}
      has_icon={false}
      onClick={on_click}
      disabled={disabled}
    >
      {label}
    </StyledStandardButton>
  );
};

export default StandardButton;

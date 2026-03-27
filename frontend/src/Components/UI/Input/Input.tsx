import React from "react";
import style from "./Input.module.scss";

interface InputProps {
  type?: "text" | "email" | "password" | "number";
  name?: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string;
  label?: string;
  error?: string;
  disabled?: boolean;
  icon?: React.ReactNode;
  borderColor?: "default" | "primary" | "success" | "warning" | "danger";
  color?: string;
  background?: string;
  size?: "small" | "medium" | "large";
  width?: string | number;
  height?: string | number;
  fullWidth?: boolean;
  borderRadius?: string | number;
  fontSize?: string | number;
}

function Input({
  type = "text",
  name,
  value,
  onChange,
  placeholder,
  label,
  error,
  disabled = false,
  icon,
  borderColor = "default",
  background,
  color,
  size = "medium",
  width,
  height,
  fullWidth = false,
  borderRadius,
  fontSize,
}: InputProps) {
  const inputStyle: React.CSSProperties = {
    width: width
      ? typeof width === "number"
        ? `${width}px`
        : width
      : undefined,
    height: height
      ? typeof height === "number"
        ? `${height}px`
        : height
      : undefined,
    borderRadius: borderRadius
      ? typeof borderRadius === "number"
        ? `${borderRadius}px`
        : borderRadius
      : undefined,
    fontSize: fontSize
      ? typeof fontSize === "number"
        ? `${fontSize}px`
        : fontSize
      : undefined,
    background: background,
    color: color,
  };

  return (
    <div
      className={`${style.inputContainer} ${fullWidth ? style.fullWidth : ""}`}
    >
      {label && <label className={style.label}>{label}</label>}
      <div className={style.inputWrapper} style={inputStyle}>
        {icon && <div className={style.icon}>{icon}</div>}
        <input
          type={type}
          name={name}
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          disabled={disabled}
          style={inputStyle}
          className={`
            ${style.input} 
            ${style[size]} 
            ${style[borderColor]} 
            ${icon ? style.withIcon : ""} 
            ${error ? style.hasError : ""}
          `}
        />
      </div>
      {error && <span className={style.errorMessage}>{error}</span>}
    </div>
  );
}

export default Input;

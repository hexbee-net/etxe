import styled from "styled-components";
import {
  space,
  width,
  height,
  fontSize,
  color,
  display,
  flex,
  flexDirection,
  alignItems,
  alignContent,
  alignSelf,
  justifyContent,
  borders,
  position,
  zIndex,
  top,
  right,
  bottom,
  left
} from "styled-system";

import { pointerEvents } from '../styles/customProperties';

export const Box = styled.div`
  ${space}
  ${width}
  ${height}
  ${fontSize}
  ${color}
  ${display}
  ${flex}
  ${flexDirection}
  ${alignContent}
  ${alignItems}
  ${alignSelf}
  ${justifyContent}
  ${borders}
  ${position}
  ${zIndex}
  ${top}
  ${right}
  ${bottom}
  ${left}
  ${pointerEvents}
`;

Box.displayName = "Box";


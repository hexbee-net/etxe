import styled from "styled-components";
import { animated } from "react-spring";

interface DrawerProps {
  width: string;
}

export const Drawer = styled(animated.div)<DrawerProps>`
  position: absolute;
  pointer-events: all;
  background-color: white;
  height: 100%;
  width: ${({ width }) => width}px;
  z-index: ${1};
`;

/** @jsxImportSource @emotion/react */

import type React from "react";
import type { FC } from 'react'
import { css } from '@emotion/react'
import {animated} from "react-spring";


const drawerStyle = (width: string, zIndex?: string) => css({
  position: "absolute",
  pointerEvents: "all",
  backgroundColor: "white",
  height: "100%",
  width: `${width}px`,
  zIndex: zIndex ?? -1,
})

interface DrawerProps {
  width: string;
  zIndex?: string;
  children?: React.ReactNode;
}

export const Drawer: FC<DrawerProps> = ({width, zIndex, children}) => (
  <animated.div css={drawerStyle(width, zIndex)}>
    {children}
  </animated.div>
)

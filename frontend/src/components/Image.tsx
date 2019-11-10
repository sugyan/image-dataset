import React, { useEffect, useState, useRef } from "react";
import { useHistory, useLocation, useParams } from "react-router";
import { Link as RouterLink, LinkProps } from "react-router-dom";
import { GlobalHotKeys } from "react-hotkeys";
import {
    Box, Grid, Link, Typography, Breadcrumbs,
    Table, TableBody, TableRow, TableCell,
} from "@material-ui/core";

import { ImageResponse } from "../common/interfaces";

const bufferLength = 100;
const bufferThreshold = 20;

const InfoTable: React.FC<ImageResponse> = (image: ImageResponse) => {
    const meta = Object.entries(JSON.parse(image.meta)).map((value, index) => {
        return (
          <Box key={index}fontFamily="Monospace">{value[0]}: {value[1]}</Box>
        );
    });
    const link = React.forwardRef<HTMLAnchorElement, Omit<LinkProps, "to">>(
        (props, ref) => {
            const params = new URLSearchParams({ name: image.label_name });
            const to = {
                pathname: "/images",
                search: params.toString(),
            };
            return (
              <RouterLink innerRef={ref} to={to} {...props} />
            );
        },
    );
    return (
      <Table>
        <TableBody>
          <TableRow>
            <TableCell component="th" scope="row">ID</TableCell>
            <TableCell><Box fontSize="h6.fontSize" fontFamily="Monospace">{image.id}</Box></TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Name</TableCell>
            <TableCell>
              <Link component={link}>{image.label_name}</Link>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Size</TableCell>
            <TableCell><Box fontSize="body1.fontSize" fontFamily="Monospace">{image.size}</Box></TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Source URL</TableCell>
            <TableCell>
              <Link href={image.source_url} target="_blank" rel="noopener">{image.source_url}</Link>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Posted at</TableCell>
            <TableCell>
              {new Date(image.posted_at * 1000).toISOString()}
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Photo URL</TableCell>
            <TableCell>
              <Link href={image.photo_url} target="_blank" rel="noopener">{image.photo_url}</Link>
            </TableCell>
          </TableRow>
          <TableRow>
            <TableCell component="th" scope="row">Meta</TableCell>
            <TableCell>
              {meta}
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    );
};

const ImageViewer: React.FC = () => {
    const history = useHistory();
    const location = useLocation();
    const params = useParams<{ id: string }>();
    const [images, setImages] = useState<ImageResponse[]>([]);
    const [terminated, setTerminated] = useState<[boolean, boolean]>([false, false]);
    const keyMap = {
        NEXT_IMAGE: ["ctrl+f", "right"],
        PREV_IMAGE: ["ctrl+b", "left"],
    };
    const nextImage = () => {
        const index = images.findIndex((element: ImageResponse) => element.id === params.id);
        if (index + 1 < images.length) {
            history.replace({
                pathname: `/image/${images[index + 1].id}`,
                search: location.search,
            });
        }
    };
    const prevImage = () => {
        const index = images.findIndex((element: ImageResponse) => element.id === params.id);
        if (index - 1 >= 0) {
            history.replace({
                pathname: `/image/${images[index - 1].id}`,
                search: location.search,
            });
        }
    };
    const handlers = {
        NEXT_IMAGE: nextImage,
        PREV_IMAGE: prevImage,
    };
    const current = images.find((element: ImageResponse) => element.id === params.id);
    const canvas = useRef<HTMLCanvasElement>(null);
    useEffect(() => {
        if (!current) {
            return;
        }
        if (canvas.current) {
            const ctx = canvas.current.getContext("2d")!;
            const image = new Image();
            image.onload = () => {
                const scale = current.size / 512;
                ctx.drawImage(image, 0, 0, 512, 512);
                ctx.strokeStyle = "cyan";
                ctx.lineWidth = 2;
                for (let i = 0; i < 68; i++) {
                    let [x, y] = [current.parts[i * 2], current.parts[i * 2 + 1]];
                    x /= scale;
                    y /= scale;
                    ctx.beginPath();
                    ctx.arc(x, y, 3, 0, 2 * Math.PI);
                    ctx.stroke();
                }
            };
            image.src = current.image_url;
        }
    }, [current]);
    useEffect(() => {
        const fetchData = async (id: string, reverse: boolean = false): Promise<ImageResponse[]> => {
            const params: URLSearchParams = new URLSearchParams(location.search);
            params.set("id", id);
            if (reverse) {
                params.set("reverse", "true");
            }
            const res = await fetch(`/api/images?${params.toString()}`);
            if (res.ok) {
                const images: ImageResponse[] = await res.json();
                if (reverse) {
                    images.reverse();
                }
                return Promise.resolve(images);
            } else {
                return Promise.reject(res.status);
            }
        };

        const requests: [Promise<ImageResponse[]>, Promise<ImageResponse[]>] = [
            Promise.resolve([]),
            Promise.resolve([]),
        ];
        if (images.length > 0) {
            const index = images.findIndex((element: ImageResponse) => element.id === params.id);
            if (index < bufferThreshold && !terminated[0]) {
                requests[0] = fetchData(images[0].id, true);
            }
            if (images.length - index <= bufferThreshold && !terminated[1]) {
                requests[1] = fetchData(images[images.length - 1].id, false);
            }
        } else {
            requests[0] = fetchData(params.id, true);
            requests[1] = fetchData(params.id, false);
        }
        Promise.all(requests).then((results: ImageResponse[][]) => {
            if (results[0].length > 0 || results[1].length > 0) {
                const map = new Map();
                [results[0], images, results[1]].flat().forEach((value: ImageResponse) => {
                    map.set(value.id, value);
                });
                const values = Array.from(map.values());
                if (values.length === images.length) {
                    if (results[0].length > 0) {
                        setTerminated([true, terminated[1]]);
                    } else {
                        setTerminated([terminated[0], true]);
                    }
                } else {
                    const indexOld = images.findIndex((element: ImageResponse) => element.id === params.id);
                    const indexNew = values.findIndex((element: ImageResponse) => element.id === params.id);
                    if (indexOld === indexNew) {
                        while (values.length > bufferLength && values[bufferThreshold].id !== params.id) {
                            values.shift();
                        }
                    } else {
                        while (values.length > bufferLength && values[values.length - bufferThreshold - 1].id !== params.id) {
                            values.pop();
                        }
                    }
                }
                setImages(values);
            }
        }).catch((err) => {
            if (err === 401) {
                history.push("/");
            } else {
                window.console.error(err);
            }
        });
    }, [location, history, params.id, images, terminated]);
    const link = React.forwardRef<HTMLAnchorElement, Omit<LinkProps, "to">>(
        (props, ref) => {
            const to = {
                pathname: "/images",
                search: location.search,
            };
            return (
              <RouterLink innerRef={ref} to={to} {...props} />
            );
        },
    );
    return (
      <div>
        <Box my={2}>
          <Breadcrumbs aria-label="breadcrumb">
            <Link color="inherit" component={link}>
              Images
            </Link>
            <Typography color="textPrimary">Image</Typography>
          </Breadcrumbs>
        </Box>
        <Grid container>
          <Grid item xs={6}>
            <canvas height={512} width={512} ref={canvas} />
          </Grid>
          <Grid>{current && <InfoTable {...current} />}</Grid>
        </Grid>
        <GlobalHotKeys keyMap={keyMap} handlers={handlers} allowChanges={true} />
      </div>
    );
};

export default ImageViewer;


const React = require('react');

const TextField = (props) => <input {...props} />;
const Button = (props) => <button {...props} />;
const Box = (props) => <div {...props} />;
const Typography = (props) => <div {...props} />;

module.exports = {
  TextField,
  Button,
  Box,
  Typography,
};

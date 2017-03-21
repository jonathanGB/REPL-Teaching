FROM msaraiva/elixir
COPY run.sh /runs/
CMD ["/runs/run.sh", "exs"]

-- step 02: seed snippets tables
insert into snippets (title, content, created, expires)
values (
		'An old silent pond',
		E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
		now(),
		now() + interval '365 days'),
	(
    'Over the wintry forest',
    E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    now(),
    now() + interval '365 days'
	),
	(
    'First autumn morning',
    E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    now(),
    now() + interval '7 days'
	);


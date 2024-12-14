-- Insert races
INSERT INTO races (name, description)
VALUES 
    ('Human', 'Guided by inner serenity and strength, Humans manifest as devout protectors of their homeland through their spiritual nature'),
    ('Numah', 'Ambitious pursuers of technology, Numahs abide to code is law and their ideals can never be overridden even if it means death')
ON CONFLICT (name) DO NOTHING;

-- Get race IDs and insert classes
WITH race_ids AS (
    SELECT id, name FROM races
)
INSERT INTO classes (race_id, name, title, description, combat_type, damage, defense, difficulty, speed)
VALUES 
    ((SELECT id FROM race_ids WHERE name = 'Human'), 'PAHLAWAN', 'The Clayborns of Power', 'A natural leader on the battlefield, the Pahlawans know what needs to be done to ensure victory prevails. Never intimidated but always intimidating, they are the symbol of motivation that drives everyone''s spirit even though they are naturally overcommitted and direct.', 'MELEE', 80, 70, 3, 50),
    ((SELECT id FROM race_ids WHERE name = 'Human'), 'PENDEKAR', 'The Guardians of Spirit', 'Masters of spiritual combat, Pendekars channel their inner strength through disciplined martial arts. Their unwavering devotion to protecting the sacred traditions of humanity makes them formidable defenders of their people''s heritage.', 'MELEE', 70, 85, 2, 45),
    ((SELECT id FROM race_ids WHERE name = 'Human'), 'PEMANAH', 'The Seekers of Truth', 'With eyes sharp as eagles and hearts pure as morning dew, Pemahs are the spiritual archers who pierce both flesh and falsehood. They maintain harmony between the physical and spiritual realms through their precise and mindful combat approach.', 'RANGED', 75, 55, 4, 65),
    ((SELECT id FROM race_ids WHERE name = 'Human'), 'PENGAMAL', 'The Mystic Sages', 'Wielders of ancient spiritual magic, Pengamals bridge the gap between the mortal and divine. Their deep understanding of life''s mysteries allows them to channel powerful magical forces while maintaining balance in the natural order.', 'MAGIC', 85, 45, 5, 40),
    ((SELECT id FROM race_ids WHERE name = 'Numah'), 'KSHATRIYA', 'The Machinists of Power', 'Built in a manner with menacing weapons, the Kshatriyas are first in line for protecting Numahs and their codes of conduct. Their duty is to never back down or refuse a challenge when their values are at stake.', 'MELEE', 85, 75, 4, 55),
    ((SELECT id FROM race_ids WHERE name = 'Numah'), 'VYAPARI', 'The Code Merchants', 'Masters of technological trade and tactical support, Vyaparis understand that true power lies in the exchange of knowledge. Their advanced support systems and resource management capabilities make them invaluable allies in any confrontation.', 'SUPPORT', 45, 60, 2, 70),
    ((SELECT id FROM race_ids WHERE name = 'Numah'), 'RAKSHAK', 'The Protocol Guardians', 'Encased in advanced defensive systems, Rakshaks are the living firewalls of Numah society. Their impenetrable defense protocols and steadfast dedication to protecting their technological sovereignty make them formidable tanks on the battlefield.', 'TANK', 60, 90, 3, 40),
    ((SELECT id FROM race_ids WHERE name = 'Numah'), 'VAIDYA', 'The System Healers', 'Specialists in both technological and biological restoration, Vaidyas maintain the perfect harmony between machine and life. Their advanced healing algorithms and support systems keep their allies operating at peak efficiency.', 'SUPPORT', 50, 65, 4, 60)
ON CONFLICT (name) DO NOTHING;

-- Function to generate random string
CREATE OR REPLACE FUNCTION random_string(length integer) RETURNS text AS $$
DECLARE
  chars text[] := '{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z}';
  result text := '';
  i integer := 0;
BEGIN
  IF length < 0 THEN
    RAISE EXCEPTION 'Given length cannot be less than 0';
  END IF;
  FOR i IN 1..length LOOP
    result := result || chars[1+random()*(array_length(chars, 1)-1)];
  END LOOP;
  RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create arrays of common username parts
DO $$
DECLARE
   prefixes text[] := ARRAY[
    'Hang', 'Puteri', 'Putera', 'Tun', 'Datuk', 'Sri', 'Raja', 'Panglima', 'Laksamana', 'Bendahara',
    'Harimau', 'Kenyalang', 'Rajawali', 'Cempaka', 'Bunga', 'Tanjung', 'Melati', 'Bayu', 'Langit',
    'Meranti', 'Rimba', 'Batik', 'Bersatu', 'Semangat', 'Keris', 'Bunga', 'Teratai', 'Sakti',
    'Badang', 'Mahsuri', 'Gunung', 'Kencana', 'Sri', 'Hikmat', 'Gemilang', 'Putih', 'Hitam', 'Perwira',
    'Pahlawan', 'Pendekar', 'Seri', 'Gagah', 'Jebat', 'Tuah', 'Melur', 'Sakti', 'Pertiwi', 'Angkasa',
    'Taufan', 'Kilauan', 'Gelanggang', 'Pelangi', 'Dewata', 'Besar', 'Tiga', 'Timur', 'Sejati', 
    'Suci', 'Damai', 'Jiwa', 'Tunas', 'Baja', 'Perkasa', 'Harmoni', 'Merah', 'Biru', 'Kuning',
    'Hijau', 'Kelapa', 'Tualang', 'Sawit', 'Perdana', 'Wawasan', 'Negaraku', 'Nadi', 'Gunung', 
    'Kinabalu', 'Petronas', 'Langkawi', 'Keris', 'Cendana', 'Sepang', 'Borneo', 'Malim', 'Sinar',
    'Mahkota', 'Cahaya', 'Selatan', 'Utara', 'Tenggara', 'Perak', 'Johor', 'Melaka', 'Sukan',
    'Hikmah', 'Andaman', 'Telaga', 'Samudera', 'Selasih', 'Kapas', 'Serumpun', 'Seri', 'Kebaya',
    'Tembok', 'Petaling', 'Rawang', 'Jati', 'Pinang', 'Layang', 'Rantau', 'Manis', 'Ampang',
     'Sultan', 'Tunku', 'Dato', 'Encik', 'Cik', 'Tok', 'Haji', 'Hajjah', 'Nik', 'Kak',
    'Abang', 'Adik', 'Pak', 'Mak', 'Lela', 'Maharaja', 'Kerajaan', 'Keraton', 'Kampung',
    'Kota', 'Bukit', 'Sungai', 'Pantai', 'Pulau', 'Lembah', 'Bujang', 'Kuil', 'Masjid',
    'Pusat', 'Tasik', 'Air', 'Api', 'Bumi', 'Lang', 'Gelora', 'Murni', 'Bersih',
    'Agung', 'Indah', 'Sari', 'Raya', 'Sukan', 'Sinar', 'Cahaya', 'Petaling', 'Rawang', 'Tualang'
];

   suffixes text[] := ARRAY[
    'Pahlawan', 'Pendekar', 'Wira', 'Jaguh', 'Perkasa', 'Harimau', 'Rajawali', 'Kenyalang', 'Mahsuri',
    'Satria', 'Helang', 'Puteri', 'Putera', 'Penyelamat', 'Langit', 'Gunung', 'Kencana', 'Gemilang',
    'Terbilang', 'Taufan', 'Bayu', 'Pelangi', 'Bunga', 'Batik', 'Semangat', 'Tuah', 'Jebat', 'Lekiu',
    'Lekir', 'Seri', 'Rimbun', 'Tangkis', 'Dewata', 'Awan', 'Bentara', 'Duta', 'Panglima', 'Laksamana',
    'Bendahara', 'Perwira', 'Pendita', 'Harmoni', 'Bersinar', 'Putih', 'Hitam', 'Bersatu', 'Melati',
    'Pertiwi', 'Jiwa', 'Biru', 'Merah', 'Hijau', 'Kuning', 'Andaman', 'Mahkota', 'Sakti', 'Damai',
    'Cendana', 'Samudera', 'Serumpun', 'Negaraku', 'Petronas', 'Sepang', 'Kinabalu', 'Langkawi', 
    'Melaka', 'Johor', 'Selasih', 'Rantau', 'Tembok', 'Manis', 'Bersih', 'Murni', 'Baja', 'Tiga',
    'Timur', 'Utara', 'Selatan', 'Tenggara', 'Tuan', 'Hamba', 'Adiwira', 'Pendita', 'Sakti', 'Bumi',
    'Rimba', 'Bendang', 'Sawah', 'Petaling', 'Rawang', 'Tualang', 'Cahaya', 'Sinar', 'Sukan', 'Besar',
    'Agung', 'Indah', 'Sari', 'Raya', 'Murni', 'Bersih', 'Gelora', 'Lang', 'Bumi', 'Api',
    'Air', 'Tasik', 'Pusat', 'Masjid', 'Kuil', 'Bujang', 'Lembah', 'Pulau', 'Pantai',
    'Sungai', 'Bukit', 'Kota', 'Kampung', 'Keraton', 'Kerajaan', 'Maharaja', 'Lela', 'Mak',
    'Pak', 'Adik', 'Abang', 'Kak', 'Nik', 'Hajjah', 'Haji', 'Tok', 'Cik', 'Encik'
    'Dato', 'Tunku', 'Sultan'
];
    names text[] := ARRAY[
    'Adam', 'Aiman', 'Farah', 'Amira', 'Afiq', 'Aina', 'Hafiz', 'Hakim', 'Nadia', 'Azizah',
    'Tuah', 'Jebat', 'Lekiu', 'Mahsuri', 'Melur', 'Seri', 'Bunga', 'Badang', 'Bayu', 'Budi',
    'Sufi', 'Putera', 'Perwira', 'Perkasa', 'Taufik', 'Iskandar', 'Syafiq', 'Alia', 'Siti', 'Nurul',
    'Zahid', 'Aisyah', 'Najib', 'Rahim', 'Salmah', 'Fazilah', 'Kamariah', 'Ali', 'Fatimah', 'Imran',
    'Azhar', 'Rania', 'Shafiqah', 'Ridwan', 'Sufian', 'Shahrul', 'Maznah', 'Latiff', 'Zainal', 'Hilmi',
    'Rosli', 'Rosnah', 'Mazlan', 'Zarina', 'Arif', 'Hamzah', 'Ahmad', 'Zaharah', 'Zulaikha', 'Hanif',
    'Kamal', 'Shafiq', 'Hidayah', 'Faizal', 'Zulkifli', 'Yasmin', 'Azman', 'Hassan', 'Shahira', 'Nazirah',
    'Ridhwan', 'Izzah', 'Zahidah', 'Nabilah', 'Shahriman', 'Nazri', 'Nur', 'Hafizah', 'Hasnah', 'Salleh',
    'Zarinah', 'Zainab', 'Zulkefli', 'Zulkarnain', 'Zul', 'Zaim', 'Zafirah', 'Zahir', 'Zaid', 'Zakaria',
    'Zarina', 'Zulaiha', 'Zulaika', 'Zubair', 'Zubaidah', 'Zulkifli', 'Zulfa', 'Zulfadli', 'Zulfiqar', 'Zulham',
    'Zulkarnain', 'Zulkhair', 'Zulkipli', 'Zulqarnain', 'Zulrafiq', 'Zulrijal', 'Zulrizal', 'Zulrukh', 'Zulsyafiq', 'Zulzakri'
];
    generated_usernames text[] := ARRAY[]::text[];
    username text;
    email text;
    password_hash text := '$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk7gpdsFGUjaP/.PeqE6nqwa'; -- 'password123'
    new_acc_id INTEGER;
    class_id INTEGER;
    new_char_id INTEGER;
    i INTEGER;
    name_style INTEGER;
BEGIN
    FOR i IN 1..50000 LOOP
        LOOP
            -- Generate username using different styles
            name_style := floor(random() * 5 + 1);
            CASE name_style
                WHEN 1 THEN -- PrefixSuffix
                    username := prefixes[floor(random() * array_length(prefixes, 1) + 1)] || 
                               suffixes[floor(random() * array_length(suffixes, 1) + 1)];
                WHEN 2 THEN -- NameNumber
                    username := names[floor(random() * array_length(names, 1) + 1)] || 
                               floor(random() * 1000)::text;
                WHEN 3 THEN -- Name_Number
                    username := names[floor(random() * array_length(names, 1) + 1)] || '_' || 
                               floor(random() * 1000)::text;
                WHEN 4 THEN -- PrefixName
                    username := prefixes[floor(random() * array_length(prefixes, 1) + 1)] || 
                               names[floor(random() * array_length(names, 1) + 1)];
                ELSE -- NameSuffix
                    username := names[floor(random() * array_length(names, 1) + 1)] || 
                               suffixes[floor(random() * array_length(suffixes, 1) + 1)];
            END CASE;

            -- Add a random string or number to ensure uniqueness
            username := username || random_string(3); -- Append a 3-character random string

            -- Check for duplicates in the array
            IF NOT (username = ANY (generated_usernames)) THEN
                -- Username is unique, add to generated list
                generated_usernames := array_append(generated_usernames, username);
                EXIT; -- Exit the loop once a unique username is generated
            END IF;
        END LOOP;

        -- Generate a unique email address
        email := lower(username) || '.' || random_string(3) || '@wira.com';

        -- Insert account
        INSERT INTO accounts (username, email, password_hash)
        VALUES (username, email, password_hash)
        RETURNING acc_id INTO new_acc_id;

        -- Create 1-3 characters for each account
        FOR j IN 1..floor(random() * 3 + 1) LOOP
            -- Get random class_id
            SELECT id INTO class_id
            FROM classes
            OFFSET floor(random() * (SELECT COUNT(*) FROM classes))
            LIMIT 1;

            -- Create character
            INSERT INTO characters (acc_id, class_id)
            VALUES (new_acc_id, class_id)
            RETURNING char_id INTO new_char_id;

            -- Add random score
            INSERT INTO scores (char_id, reward_score)
            VALUES (new_char_id, floor(random() * 10000));
        END LOOP;
    END LOOP;
END $$;
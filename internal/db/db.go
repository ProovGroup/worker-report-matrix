package db

import (
	"encoding/json"
	"fmt"
	r "worker-report-matrix/internal/report"

	"github.com/ProovGroup/env"
)

const QUERY = `SELECT row_to_json(matrix)
		FROM (
		SELECT
			r.proov_code AS proov_code,
			get_local_or_en(r.title, 'fr') AS title,
			i.value AS identifier_item,
			'data:image/png;base64,' || c.data as logo,
			r."createdAt" AS created_at,
			r.owner AS owner,
			CAST(r.geoloc->>'latitude' as FLOAT8) AS latitude,
			CAST(r.geoloc->>'longitude' as FLOAT8) AS longitude,
			(SELECT row_to_json(cag)
				FROM (
					SELECT
						cag.latitude,
						cag.longitude,
						cag.address,
						'' AS timezone
					FROM
						cache_address_geoloc cag
					WHERE 
						cag.latitude = CAST(r.geoloc->>'latitude' AS float8) and cag.longitude = CAST(r.geoloc->>'longitude' AS float8)
				) AS cag
			) AS geoloc,
			(SELECT COALESCE(array_to_json(array_agg(row_to_json(parts))), '[]') 
				FROM (
					SELECT 
						get_local_or_en(pt.title, 'fr') AS title,
						pt.note,
						(SELECT s3.path FROM s3picture s3 WHERE s3.id = pt.sign) AS sign_url,
						(SELECT COALESCE(array_to_json(array_agg(row_to_json(iss))) , '[]')
							FROM (SELECT 
								get_local_or_en(i.title, 'fr') AS title,
								i.value
							FROM rptinfos i
							WHERE i.part_id = pt.id ORDER BY i.position) AS iss
						) AS infos
					FROM rptpart pt 
					WHERE pt.report_id = r.id
					ORDER BY pt.position
				) AS parts
			) AS parts,
			(SELECT COALESCE(array_to_json(array_agg(row_to_json(ps))), '[]') 
				FROM (
					SELECT 
						ps.name,
						(SELECT row_to_json(pic)
							FROM (
								SELECT
									s3.path AS original,
									SUBSTRING(s3.path, 0, POSITION('.' IN s3.path)) || '/thumb.jpeg' AS thumbnail
								FROM s3picture s3 
								WHERE s3.id = ps.picture_id
							) AS pic
						) AS picture,						
						(SELECT array_to_json(array_agg(row_to_json(d)))
							FROM (
								SELECT
									d.matrix_price,
									d.matrix_entries AS matrix_answers
								FROM
									rptdamages d
								WHERE
									d.rptprocess_id = ps.id AND d.is_drop_off = false
							) AS d
						) AS infos_damages
					FROM 
						rptprocess ps 
					WHERE ps.report_id = r.id
					ORDER BY ps.position
				) AS ps
			) AS processes
		FROM
			report r
		LEFT JOIN items i ON r.item_id = i.id
		LEFT JOIN cachedimg c ON r.image_report = c.id
		WHERE
			proov_code = $1
		) AS matrix`

func GetReport(e *env.Env, proovCode string) (*r.Report, error) {
	var data []byte
	err := e.QueryRow(QUERY, proovCode).Scan(&data)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("Empty result set for proov_code: %s", proovCode)
	}

	var report r.Report
	err = json.Unmarshal(data, &report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func SaveGeoLoc(e *env.Env, geoloc *r.GeoLoc) error {
	query := `INSERT INTO cache_address_geoloc (latitude, longitude, locale, address, created_at)
		VALUES ($1, $2, 'fr', $3, now())`

	_, err := e.Exec(query, geoloc.Latitude, geoloc.Longitude, geoloc.Address)
	return err
}
